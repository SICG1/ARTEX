import packageJson from "../../package.json";

import { t } from "@/i18n";

const currentYear = new Date().getFullYear();

export const APP_CONFIG = {
  name: "Security Testing Center",
  version: packageJson.version,
  copyright: `© ${currentYear}, Security Testing Center.`,
  meta: {
    title: t("ARTEX — Autonomous Penetration Testing Console"),
    description: t("Console for the LLM-driven autonomous penetration testing system"),
  },
};

