#!/usr/bin/env npx tsx
import type { CliOptions } from "../core/types.js";
import { applyProvider, listProviders, validateConfig } from "../services/validate-service.js";

function printHelp(): void {
  console.log("Usage:");
  console.log("  npx tsx src/cli/index.ts [provider-id|provider-name] [--dry-run]");
  console.log("  npx tsx src/cli/index.ts --list");
  console.log("  npx tsx src/cli/index.ts --validate");
  console.log("  npx tsx src/cli/index.ts --help");
  console.log("");
  console.log("Options:");
  console.log("  --list           List providers (default when no args)");
  console.log("  --validate       Validate providers.yaml and run live checks");
  console.log("  --dry-run        Preview switch plan without writing files");
  console.log("  --help, -h       Show help");
  console.log("");
  console.log("Examples:");
  console.log("  npx tsx src/cli/index.ts");
  console.log("  npx tsx src/cli/index.ts openrouter");
  console.log("  npx tsx src/cli/index.ts Codex2Claude CLIProxy API");
  console.log("  npx tsx src/cli/index.ts openrouter --dry-run");
  console.log("  npx tsx src/cli/index.ts --validate");
}

function parseCliArgs(argv: string[]): CliOptions {
  let dryRun = false;
  let action: CliOptions["action"] | null = null;
  const positional: string[] = [];

  for (const arg of argv) {
    if (arg === "--dry-run") {
      dryRun = true;
      continue;
    }

    if (arg === "--help" || arg === "-h") {
      if (action && action !== "help") {
        throw new Error("--help cannot be combined with other actions");
      }
      action = "help";
      continue;
    }

    if (arg === "--list") {
      if (action && action !== "list") {
        throw new Error("--list cannot be combined with other actions");
      }
      action = "list";
      continue;
    }

    if (arg === "--validate" || arg === "validate") {
      if (action && action !== "validate") {
        throw new Error("--validate cannot be combined with other actions");
      }
      action = "validate";
      continue;
    }

    if (arg.startsWith("-")) {
      throw new Error(`Unknown option: ${arg}`);
    }

    positional.push(arg);
  }

  if (action === "help") {
    if (dryRun) {
      throw new Error("--dry-run cannot be used with --help");
    }
    if (positional.length > 0) {
      throw new Error("Unexpected provider argument with --help");
    }
    return { action: "help", dryRun: false };
  }

  if (action === "list") {
    if (dryRun) {
      throw new Error("--dry-run can only be used when switching provider");
    }
    if (positional.length > 0) {
      throw new Error("Unexpected provider argument with --list");
    }
    return { action: "list", dryRun: false };
  }

  if (action === "validate") {
    if (dryRun) {
      throw new Error("--dry-run can only be used when switching provider");
    }
    if (positional.length > 0) {
      throw new Error("Unexpected provider argument with --validate");
    }
    return { action: "validate", dryRun: false };
  }

  if (positional.length === 0) {
    return { action: "list", dryRun: false };
  }

  if (positional.length > 1) {
    throw new Error(`Too many provider arguments: ${positional.join(", ")}`);
  }

  return {
    action: "switch",
    providerInput: positional[0],
    dryRun,
  };
}

try {
  const cli = parseCliArgs(process.argv.slice(2));

  if (cli.action === "help") {
    printHelp();
  } else if (cli.action === "list") {
    listProviders();
  } else if (cli.action === "validate") {
    validateConfig();
  } else {
    applyProvider(cli.providerInput!, { dryRun: cli.dryRun });
  }
} catch (error) {
  const message = error instanceof Error ? error.message : String(error);
  console.error(`Error: ${message}`);
  console.error("Tip: run with --help to see supported usage.");
  process.exit(1);
}
