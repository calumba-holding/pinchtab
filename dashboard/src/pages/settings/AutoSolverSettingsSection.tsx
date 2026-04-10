import type { BackendConfig, BackendConfigState } from "../../types";
import type { UpdateBackendSection } from "./settingsShared";
import { csvToList, fieldClass, listToCsv } from "./settingsShared";
import { SectionCard, SettingRow } from "./SettingsSharedComponents";

interface AutoSolverSettingsSectionProps {
  backendConfig: BackendConfig;
  backendState: BackendConfigState | null;
  updateBackendSection: UpdateBackendSection;
}

export function AutoSolverSettingsSection({
  backendConfig,
  backendState,
  updateBackendSection,
}: AutoSolverSettingsSectionProps) {
  return (
    <SectionCard
      title="AutoSolver"
      description="These settings are saved into the PinchTab config file. External provider API keys stay write-only and must be set directly in that file."
    >
      <SettingRow
        label="Config file"
        description="Dashboard edits are written back to this file. Set external provider keys under autoSolver.external in the same config file."
      >
        <div className="rounded-sm border border-border-subtle bg-[rgb(var(--brand-surface-code-rgb)/0.72)] px-3 py-2 text-sm text-text-secondary">
          <code>{backendState?.configPath || "Config path unavailable"}</code>
        </div>
      </SettingRow>
      <SettingRow
        label="Enable AutoSolver"
        description="Turns on the autosolver runtime configuration for supported challenge flows."
      >
        <label className="flex items-center justify-end gap-3 text-sm text-text-secondary">
          <input
            type="checkbox"
            checked={backendConfig.autoSolver.enabled}
            onChange={(e) =>
              updateBackendSection("autoSolver", {
                enabled: e.target.checked,
              })
            }
            className="h-4 w-4"
          />
          {backendConfig.autoSolver.enabled ? "Enabled" : "Disabled"}
        </label>
      </SettingRow>
      <SettingRow
        label="Max attempts"
        description="Maximum autosolver retries before the pipeline gives up."
      >
        <input
          type="number"
          min={1}
          value={backendConfig.autoSolver.maxAttempts}
          onChange={(e) =>
            updateBackendSection("autoSolver", {
              maxAttempts: Number(e.target.value),
            })
          }
          className={fieldClass}
        />
      </SettingRow>
      <SettingRow
        label="Solvers"
        description="Comma-separated ordered list of solver names to try."
      >
        <input
          value={listToCsv(backendConfig.autoSolver.solvers)}
          onChange={(e) =>
            updateBackendSection("autoSolver", {
              solvers: csvToList(e.target.value),
            })
          }
          className={fieldClass}
        />
      </SettingRow>
      <SettingRow
        label="LLM provider"
        description="Optional provider name used when LLM fallback is enabled."
      >
        <input
          value={backendConfig.autoSolver.llmProvider}
          onChange={(e) =>
            updateBackendSection("autoSolver", {
              llmProvider: e.target.value,
            })
          }
          className={fieldClass}
        />
      </SettingRow>
      <SettingRow
        label="LLM fallback"
        description="Use an LLM as the last resort after registered solvers fail."
      >
        <label className="flex items-center justify-end gap-3 text-sm text-text-secondary">
          <input
            type="checkbox"
            checked={backendConfig.autoSolver.llmFallback}
            onChange={(e) =>
              updateBackendSection("autoSolver", {
                llmFallback: e.target.checked,
              })
            }
            className="h-4 w-4"
          />
          {backendConfig.autoSolver.llmFallback ? "Enabled" : "Disabled"}
        </label>
      </SettingRow>
      <SettingRow
        label="External provider keys"
        description="Capsolver and 2Captcha credentials are not shown in the dashboard and must be managed in the config file."
      >
        <div className="rounded-sm border border-warning/25 bg-warning/10 px-3 py-2 text-xs leading-5 text-warning">
          Open the config file above and set{" "}
          <code>autoSolver.external.capsolverKey</code> and{" "}
          <code>autoSolver.external.twoCaptchaKey</code> there. The dashboard
          does not display or edit those values, and there are no environment
          variable overrides.
        </div>
      </SettingRow>
    </SectionCard>
  );
}
