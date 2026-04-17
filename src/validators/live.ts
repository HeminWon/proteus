import { spawnSync } from "child_process";

import type { Provider } from "./structure.js";

export interface LiveValidationResult {
  providerId: string;
  status: "ok" | "fail" | "skip";
  detail: string;
  latencyMs: number | null;
}

function trimTrailingSlash(value: string): string {
  return value.replace(/\/+$/, "");
}

function buildModelsUrl(baseUrl: string): string {
  const base = trimTrailingSlash(baseUrl);
  if (base.includes("openrouter.ai/api/v1")) {
    return `${base}/auth/key`;
  }
  if (base.endsWith("/v1")) {
    return `${base}/models`;
  }
  return `${base}/v1/models`;
}

function maskToken(token: string): string {
  if (token.length <= 8) {
    return "****";
  }
  return `${token.slice(0, 4)}...${token.slice(-4)}`;
}

export function validateProviderLive(provider: Provider): LiveValidationResult {
  const startedAt = Date.now();
  const env = provider.claude.env;
  const token = env.ANTHROPIC_AUTH_TOKEN;
  const baseUrl = env.ANTHROPIC_BASE_URL;

  if (!baseUrl || typeof baseUrl !== "string") {
    return {
      providerId: provider.id,
      status: "skip",
      detail: "missing ANTHROPIC_BASE_URL",
      latencyMs: null,
    };
  }

  if (!token || typeof token !== "string") {
    return {
      providerId: provider.id,
      status: "skip",
      detail: "missing ANTHROPIC_AUTH_TOKEN",
      latencyMs: null,
    };
  }

  const url = buildModelsUrl(baseUrl);
  const args = [
    "-sS",
    "-L",
    "--max-time",
    "15",
    "-o",
    "/dev/null",
    "-w",
    "%{http_code}",
    "-H",
    `x-api-key: ${token}`,
    "-H",
    `authorization: Bearer ${token}`,
    "-H",
    "anthropic-version: 2023-06-01",
    "-H",
    "accept: application/json",
    url,
  ];

  const result = spawnSync("curl", args, {
    encoding: "utf8",
    timeout: 20000,
  });

  if (result.error) {
    return {
      providerId: provider.id,
      status: "fail",
      detail: `curl error: ${result.error.message}`,
      latencyMs: Date.now() - startedAt,
    };
  }

  if (result.status !== 0) {
    const err = (result.stderr || "").trim();
    return {
      providerId: provider.id,
      status: "fail",
      detail: `curl exited ${result.status}${err ? `: ${err}` : ""}`,
      latencyMs: Date.now() - startedAt,
    };
  }

  const code = Number((result.stdout || "").trim());
  if (!Number.isFinite(code)) {
    return {
      providerId: provider.id,
      status: "fail",
      detail: "invalid HTTP status returned by curl",
      latencyMs: Date.now() - startedAt,
    };
  }

  const latencyMs = Date.now() - startedAt;

  if (code >= 200 && code < 300) {
    return {
      providerId: provider.id,
      status: "ok",
      detail: `HTTP ${code} (${url}, token ${maskToken(token)})`,
      latencyMs,
    };
  }

  return {
    providerId: provider.id,
    status: "fail",
    detail: `HTTP ${code} (${url}, token ${maskToken(token)})`,
    latencyMs,
  };
}
