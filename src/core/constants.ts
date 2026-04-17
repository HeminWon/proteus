import * as path from "path";

export const SETTINGS_PATH = path.join(process.env.HOME!, ".claude", "settings.json");

export const CACHE_PATH = path.join(
  process.env.XDG_CACHE_HOME ?? path.join(process.env.HOME!, ".cache"),
  "proteus",
  "cache.json"
);

export const BACKUP_DIR = path.join(process.env.HOME!, ".claude", "proteus-backups");

export const MAX_BACKUPS = 10;
export const HIGH_LATENCY_MS = 1500;
