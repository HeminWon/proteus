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

export async function validateProviderLive(provider: Provider): Promise<LiveValidationResult> {
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
  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), 20000);

  try {
    const response = await fetch(url, {
      method: "GET",
      headers: {
        "x-api-key": token,
        authorization: `Bearer ${token}`,
        "anthropic-version": "2023-06-01",
        accept: "application/json",
      },
      redirect: "follow",
      signal: controller.signal,
    });

    const latencyMs = Date.now() - startedAt;
    const code = response.status;

    if (code >= 200 && code < 300) {
      return {
        providerId: provider.id,
        status: "ok",
        detail: `HTTP ${code} (${url})`,
        latencyMs,
      };
    }

    let bodySnippet = "";
    try {
      const body = (await response.text()).trim();
      if (body.length > 0) {
        bodySnippet = body.replace(/\s+/g, " ").slice(0, 200);
      }
    } catch {
      // ignore body parse errors
    }

    return {
      providerId: provider.id,
      status: "fail",
      detail: bodySnippet
        ? `HTTP ${code} (${url}) | body=${JSON.stringify(bodySnippet)}`
        : `HTTP ${code} (${url})`,
      latencyMs,
    };
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error);
    return {
      providerId: provider.id,
      status: "fail",
      detail: `request error: ${message}`,
      latencyMs: Date.now() - startedAt,
    };
  } finally {
    clearTimeout(timeout);
  }
}
