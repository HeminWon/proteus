import * as fs from "fs";
import * as path from "path";
import * as yaml from "js-yaml";

import type { ProvidersConfig } from "./types.js";
import { validateProvidersConfigShape } from "../validators/structure.js";

export interface LoadProvidersResult {
  config: ProvidersConfig;
  configDir: string;
}

function readAppConfig(): { config_dir?: string } {
  const xdgConfigHome = process.env.XDG_CONFIG_HOME ?? path.join(process.env.HOME!, ".config");
  const configPath = path.join(xdgConfigHome, "proteus", "config.json");
  if (!fs.existsSync(configPath)) return {};
  try {
    const parsed = JSON.parse(fs.readFileSync(configPath, "utf8"));
    if (parsed && typeof parsed === "object" && !Array.isArray(parsed)) {
      return parsed as { config_dir?: string };
    }
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error);
    console.warn(`Warning: failed to parse ${configPath}: ${message}`);
  }
  return {};
}

export function resolveConfigDir(): string {
  const xdgConfigHome = process.env.XDG_CONFIG_HOME ?? path.join(process.env.HOME!, ".config");
  const xdgProteusDir = path.join(xdgConfigHome, "proteus");

  const appConfig = readAppConfig();
  const candidates: Array<{ dir: string; label: string }> = [];

  if (appConfig.config_dir) {
    const expanded = appConfig.config_dir.replace(/^~/, process.env.HOME!);
    candidates.push({ dir: expanded, label: "config.json (config_dir)" });
  }
  candidates.push({ dir: xdgProteusDir, label: `XDG (${xdgProteusDir})` });

  for (const { dir } of candidates) {
    if (fs.existsSync(path.join(dir, "providers.yaml"))) return dir;
  }

  const configJsonPath = path.join(xdgProteusDir, "config.json");
  throw new Error(
    `providers.yaml not found.\nSearched:\n` +
      candidates.map((c) => `  - ${path.join(c.dir, "providers.yaml")}`).join("\n") +
      `\n\nPlace your providers.yaml at:\n  ${path.join(xdgProteusDir, "providers.yaml")}\n` +
      `Or set config_dir in ${configJsonPath}`
  );
}

export function loadProviders(): LoadProvidersResult {
  const configDir = resolveConfigDir();
  const parsed = yaml.load(fs.readFileSync(path.join(configDir, "providers.yaml"), "utf8"));
  return { config: validateProvidersConfigShape(parsed), configDir };
}
