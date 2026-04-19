import * as fs from "fs";
import * as path from "path";

import {
  BACKUP_DIR,
  CACHE_PATH,
  HIGH_LATENCY_MS,
  MAX_BACKUPS,
  SETTINGS_PATH,
} from "../core/constants.js";
import type {
  CacheData,
  JsonObject,
  SettingsReadResult,
  SwitchPlan,
} from "../core/types.js";
import { loadProviders } from "../providers/loader.js";
import type { Provider, ProvidersConfig } from "../providers/types.js";
import { validateProviderLive, type LiveValidationResult } from "../validators/live.js";

function supportsColor(): boolean {
  return Boolean(process.stdout.isTTY) && process.env.NO_COLOR === undefined;
}

function colorize(text: string, colorCode: string): string {
  if (!supportsColor()) {
    return text;
  }
  return `\u001b[${colorCode}m${text}\u001b[0m`;
}

function colorStatus(status: "ok" | "fail" | "skip", mark: string): string {
  if (status === "ok") {
    return colorize(mark, "32");
  }
  if (status === "fail") {
    return colorize(mark, "31");
  }
  return colorize(mark, "33");
}

function formatLatency(latencyMs: number | null, status: "ok" | "fail" | "skip"): string {
  if (latencyMs === null) {
    return colorize("n/a", "33");
  }

  const raw = `${latencyMs}ms`;
  if (latencyMs >= HIGH_LATENCY_MS) {
    return colorize(raw, "33");
  }

  if (status === "ok") {
    return colorize(raw, "32");
  }

  if (status === "fail") {
    return colorize(raw, "31");
  }

  return raw;
}

function findProviderById(config: ProvidersConfig, id: string): Provider | undefined {
  return config.providers.find((provider) => provider.id === id);
}

function findProviderByInput(config: ProvidersConfig, input: string): Provider | undefined {
  const byId = findProviderById(config, input);
  if (byId) {
    return byId;
  }
  return config.providers.find((provider) => provider.name === input);
}

function writeFileAtomic(filePath: string, content: string): void {
  const dir = path.dirname(filePath);
  fs.mkdirSync(dir, { recursive: true });
  const tmpPath = `${filePath}.tmp-${process.pid}-${Date.now()}`;
  fs.writeFileSync(tmpPath, content, "utf8");
  fs.renameSync(tmpPath, filePath);
}

function readCache(): CacheData {
  if (!fs.existsSync(CACHE_PATH)) {
    return {};
  }

  let parsed: unknown;
  try {
    parsed = JSON.parse(fs.readFileSync(CACHE_PATH, "utf8"));
  } catch {
    return {};
  }

  if (!parsed || typeof parsed !== "object" || Array.isArray(parsed)) {
    return {};
  }

  return parsed as CacheData;
}

function writeCache(cache: CacheData): void {
  writeFileAtomic(CACHE_PATH, `${JSON.stringify(cache, null, 2)}\n`);
}

function getActiveProviderId(config: ProvidersConfig, cache: CacheData): string | null {
  const cached = cache.active?.claude;
  if (cached && findProviderById(config, cached)) {
    return cached;
  }

  return null;
}

function readSettings(): SettingsReadResult {
  if (!fs.existsSync(SETTINGS_PATH)) {
    return { exists: false, data: {} };
  }

  const raw = fs.readFileSync(SETTINGS_PATH, "utf8");
  let parsed: unknown;
  try {
    parsed = JSON.parse(raw);
  } catch (error) {
    throw new Error(`Failed to parse ${SETTINGS_PATH}: ${(error as Error).message}`);
  }

  if (!parsed || typeof parsed !== "object" || Array.isArray(parsed)) {
    throw new Error(`Invalid settings root in ${SETTINGS_PATH}: expected JSON object`);
  }

  return { exists: true, data: parsed as JsonObject };
}

function writeSettings(settings: JsonObject): void {
  writeFileAtomic(SETTINGS_PATH, `${JSON.stringify(settings, null, 2)}\n`);
}

function formatTimestamp(date: Date): string {
  const yyyy = date.getFullYear();
  const mm = String(date.getMonth() + 1).padStart(2, "0");
  const dd = String(date.getDate()).padStart(2, "0");
  const hh = String(date.getHours()).padStart(2, "0");
  const mi = String(date.getMinutes()).padStart(2, "0");
  const ss = String(date.getSeconds()).padStart(2, "0");
  return `${yyyy}${mm}${dd}_${hh}${mi}${ss}`;
}

function createBackupIfNeeded(settingsExists: boolean): string | null {
  if (!settingsExists) {
    return null;
  }

  fs.mkdirSync(BACKUP_DIR, { recursive: true });
  const backupName = `settings-${formatTimestamp(new Date())}.json`;
  const backupPath = path.join(BACKUP_DIR, backupName);
  fs.copyFileSync(SETTINGS_PATH, backupPath);
  cleanupOldBackups();
  return backupPath;
}

function cleanupOldBackups(): void {
  if (!fs.existsSync(BACKUP_DIR)) {
    return;
  }

  const backups = fs
    .readdirSync(BACKUP_DIR)
    .filter((name) => /^settings-\d{8}_\d{6}\.json$/.test(name))
    .map((name) => {
      const fullPath = path.join(BACKUP_DIR, name);
      return {
        fullPath,
        mtimeMs: fs.statSync(fullPath).mtimeMs,
      };
    })
    .sort((a, b) => b.mtimeMs - a.mtimeMs);

  for (const stale of backups.slice(MAX_BACKUPS)) {
    fs.unlinkSync(stale.fullPath);
  }
}

function restoreFromBackup(backupPath: string): void {
  if (!backupPath || !fs.existsSync(backupPath)) {
    return;
  }
  fs.copyFileSync(backupPath, SETTINGS_PATH);
}

function asStringMap(value: unknown): Record<string, string> {
  if (!value || typeof value !== "object" || Array.isArray(value)) {
    return {};
  }

  const map: Record<string, string> = {};
  for (const [key, raw] of Object.entries(value as Record<string, unknown>)) {
    if (typeof raw === "string") {
      map[key] = raw;
    }
  }
  return map;
}

function buildNextSettings(
  _config: ProvidersConfig,
  currentSettings: JsonObject,
  provider: Provider
): JsonObject {
  const next: JsonObject = { ...currentSettings };
  next.env = { ...provider.claude.env };

  if (provider.claude.models) {
    next.availableModels = [...provider.claude.models];
  } else {
    delete next.availableModels;
  }

  return next;
}

function createSwitchPlan(
  activeProviderId: string | null,
  currentSettings: JsonObject,
  nextSettings: JsonObject,
  targetProvider: string,
  backupRequired: boolean
): SwitchPlan {
  const beforeEnv = asStringMap(currentSettings.env);
  const afterEnv = asStringMap(nextSettings.env);

  const envAdded: string[] = [];
  const envUpdated: string[] = [];
  const envRemoved: string[] = [];

  for (const key of Object.keys(afterEnv)) {
    if (!(key in beforeEnv)) {
      envAdded.push(key);
    } else if (beforeEnv[key] !== afterEnv[key]) {
      envUpdated.push(key);
    }
  }

  for (const key of Object.keys(beforeEnv)) {
    if (!(key in afterEnv)) {
      envRemoved.push(key);
    }
  }

  const beforeModels = JSON.stringify(currentSettings.availableModels ?? null);
  const afterModels = JSON.stringify(nextSettings.availableModels ?? null);

  return {
    fromProvider: activeProviderId,
    toProvider: targetProvider,
    envAdded: envAdded.sort(),
    envUpdated: envUpdated.sort(),
    envRemoved: envRemoved.sort(),
    availableModelsChanged: beforeModels !== afterModels,
    backupRequired,
  };
}

function printSwitchPlan(plan: SwitchPlan): void {
  console.log(`Plan: ${plan.fromProvider ?? "(unset)"} -> ${plan.toProvider}`);
  console.log("- mode: overwrite-env");
  console.log(`- settings: ${SETTINGS_PATH}`);
  console.log(`- cache: ${CACHE_PATH}`);
  console.log(`- backup: ${plan.backupRequired ? "create" : "skip (file missing)"}`);

  const changes = [
    ["env added", plan.envAdded],
    ["env updated", plan.envUpdated],
    ["env removed", plan.envRemoved],
  ] as const;

  for (const [label, keys] of changes) {
    console.log(`- ${label}: ${keys.length > 0 ? keys.join(", ") : "none"}`);
  }

  console.log(`- availableModels: ${plan.availableModelsChanged ? "changed" : "unchanged"}`);
}

export function applyProvider(name: string, options: { dryRun: boolean }): void {
  const { config } = loadProviders();
  const cache = readCache();
  const activeProviderId = getActiveProviderId(config, cache);
  const provider = findProviderByInput(config, name);
  if (!provider) {
    const available = config.providers.map((item) => item.id).join(", ");
    throw new Error(`Provider "${name}" not found. Available: ${available}`);
  }

  const { exists: settingsExists, data: currentSettings } = readSettings();
  const nextSettings = buildNextSettings(config, currentSettings, provider);

  const plan = createSwitchPlan(
    activeProviderId,
    currentSettings,
    nextSettings,
    provider.id,
    settingsExists
  );

  if (options.dryRun) {
    printSwitchPlan(plan);
    return;
  }

  const backupPath = createBackupIfNeeded(settingsExists);

  try {
    writeSettings(nextSettings);
    writeCache({ active: { claude: provider.id } });
  } catch (error) {
    if (backupPath) {
      try {
        restoreFromBackup(backupPath);
      } catch {
        // Best effort restore.
      }
    }
    throw error;
  }

  console.log(`Switched to: ${provider.name} (${provider.id})`);
  if (backupPath) {
    console.log(`Backup: ${backupPath}`);
  }
  console.log("Mode: overwrite-env");
}

export function listProviders(): void {
  const { config, configDir } = loadProviders();
  const cache = readCache();
  const activeProviderId = getActiveProviderId(config, cache);
  console.log(`Config dir: ${configDir}\n`);
  console.log("Available providers:\n");
  for (const provider of config.providers) {
    const active = activeProviderId && provider.id === activeProviderId ? " ◀ active" : "";
    console.log(`  ${provider.id.padEnd(16)} ${provider.name}${active}`);
  }
}

async function validateProvidersLive(
  providers: Provider[],
  concurrency: number
): Promise<LiveValidationResult[]> {
  const results: LiveValidationResult[] = [];

  for (let i = 0; i < providers.length; i += concurrency) {
    const batch = providers.slice(i, i + concurrency);
    const batchResults = await Promise.all(batch.map((provider) => validateProviderLive(provider)));
    results.push(...batchResults);
  }

  return results;
}

export async function validateConfig(): Promise<void> {
  const { config, configDir } = loadProviders();
  const cache = readCache();
  const activeProviderId = getActiveProviderId(config, cache);
  console.log("providers.yaml is valid.");
  console.log(`- config dir: ${configDir}`);
  console.log(`- version: ${config.version}`);
  console.log(`- active (cache): ${activeProviderId ?? "unset"}`);
  console.log(`- providers: ${config.providers.length}`);

  const concurrency = 5;
  console.log(`- live validation: enabled (HTTP endpoint, concurrency=${concurrency})`);

  const results = await validateProvidersLive(config.providers, concurrency);
  for (const result of results) {
    const mark = result.status === "ok" ? "OK" : result.status === "skip" ? "SKIP" : "FAIL";
    const markDisplay = colorStatus(result.status, mark);
    const latencyDisplay = formatLatency(result.latencyMs, result.status);
    console.log(`  [${markDisplay}] ${result.providerId}: ${result.detail} | latency=${latencyDisplay}`);
  }

  const failed = results.filter((result) => result.status === "fail");
  if (failed.length > 0) {
    throw new Error(`live validation failed for ${failed.length} provider(s)`);
  }
}
