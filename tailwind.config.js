import franken from "franken-ui/shadcn-ui/preset-quick";

/** @type {import('tailwindcss').Config} */
export default {
  presets: [
    franken({
      customPalette: {
        ".ux-theme-catppuccin": {
          "--background": "220 23.08% 94.9%",
          "--foreground": "233.79 16.02% 35.49%",
          "--card": "220 21.95% 91.96%",
          "--card-foreground": "233.79 16.02% 35.49%",
          "--popover": "220 20.69% 88.63%",
          "--popover-foreground": "233.79 16.02% 35.49%",
          "--primary": "219.91 91.49% 53.92%",
          "--primary-foreground": "0 0% 96.08%",
          "--secondary": "226.67 12.16% 70.98%",
          "--secondary-foreground": "240 22.73% 8.63%",
          "--muted": "225 13.56% 76.86%",
          "--muted-foreground": "240 21.05% 14.9%",
          "--accent": "230.94 97.2% 71.96%",
          "--accent-foreground": "240 22.73% 8.63%",
          "--destructive": "0 84.2% 60.2%",
          "--destructive-foreground": "210 40% 98%",
          "--border": "233.79 16.02% 35.49%",
          "--input": "226.67 12.16% 70.98%",
          "--ring": "197.07 96.57% 45.69%",
        },
        ".dark.ux-theme-catppuccin": {
          "--background": "240 21.05% 14.9%",
          "--foreground": "226.15 63.93% 88.04%",
          "--card": "0 0% 7.84%",
          "--card-foreground": "226.15 63.93% 88.04%",
          "--popover": "240 22.73% 8.63%",
          "--popover-foreground": "226.15 63.93% 88.04%",
          "--primary": "217.17 91.87% 75.88%",
          "--primary-foreground": "240 22.73% 8.63%",
          "--secondary": "232.5 12% 39.22%",
          "--secondary-foreground": "226.15 63.93% 88.04%",
          "--muted": "236.84 16.24% 22.94%",
          "--muted-foreground": "226.15 63.93% 88.04%",
          "--accent": "115.45 54.1% 76.08%",
          "--accent-foreground": "240 22.73% 8.63%",
          "--destructive": "343 81.2% 74.9%",
          "--destructive-foreground": "240 22.73% 8.63%",
          "--border": "236.84 16.24% 22.94%",
          "--input": "232.5 12% 39.22%",
          "--ring": "189.18 71.01% 72.94%",
        },
      },
    }),
  ],
  content: ["./view/**/*.templ"],
  safelist: [
    {
      pattern: /^uk-/,
    },
    "ProseMirror",
    "ProseMirror-focused",
    "tiptap",
  ],
  theme: {
    extend: {
      colors: {
        warning: "hsl(var(--warning))",
        "warning-foreground": "hsl(var(--warning-foreground))",
        success: "hsl(var(--success))",
        "success-foreground": "hsl(var(--success-foreground))",
        "chart-green": "hsl(var(--chart-green)",
        "chart-blue": "hsl(var(--chart-blue)",
        "chart-yellow": "hsl(var(--chart-yellow)",
        "chart-orange": "hsl(var(--chart-orange)",
        "chart-red": "hsl(var(--chart-red)",
      },
    },
  },
  plugins: [],
};
