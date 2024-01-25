import { Ref, computed, nextTick, ref, unref } from "vue";
import { useI18n } from "vue-i18n";
import { ComposedProject, MaybeRef } from "@/types";
import {
  ExtractPromiseType,
  extractProjectResourceName,
  minmax,
} from "@/utils";

export type UseRenderMarkdownOptions = {
  minHeight: number;
  maxHeight: number;
  placeholder: string;
};

const defaultOptions = (): UseRenderMarkdownOptions => ({
  minHeight: 48 /* 3rem */,
  maxHeight: 192 /* 12rem */,
  placeholder: "",
});

export const useRenderMarkdown = (
  markdown: MaybeRef<string>,
  iframeRef: Ref<HTMLIFrameElement | undefined>,
  projectRef: MaybeRef<ComposedProject | undefined> = ref(),
  options: Partial<UseRenderMarkdownOptions> | undefined = undefined
) => {
  const mergedOptions = Object.assign(defaultOptions(), options);
  const { t } = useI18n();
  const request = Promise.all([
    import("highlight.js/lib/core"),
    import("highlight.js/styles/github.css?raw"),
    import("@/assets/css/github-markdown-style.css?raw"),
    import("markdown-it"),
    import("dompurify"),
  ]);
  const modules = ref<ExtractPromiseType<typeof request>>();
  request.then((dep) => {
    modules.value = dep;
  });
  const deps = computed(() => {
    if (!modules.value) return undefined;
    const [
      { default: hljs },
      { default: codeStyle },
      { default: markdownStyle },
      { default: MarkdownIt },
      { default: DOMPurify },
    ] = modules.value;
    const md = new MarkdownIt({
      html: true,
      linkify: true,
      highlight: function (code, lang) {
        if (lang && hljs.getLanguage(lang)) {
          try {
            return hljs.highlight(code, { language: lang }).value;
          } catch {
            return "";
          }
        }

        return ""; // use external default escaping
      },
    });

    // See: https://github.com/cure53/DOMPurify/tree/main/demos#hook-to-open-all-links-in-a-new-window-link
    // Add a hook to make all links open a new window
    DOMPurify.addHook("afterSanitizeAttributes", function (node) {
      // set all elements owning target to target=_blank
      if ("target" in node) {
        node.setAttribute("target", "_blank");
      }
      // set non-HTML/MathML links to xlink:show=new
      if (
        !node.hasAttribute("target") &&
        (node.hasAttribute("xlink:href") || node.hasAttribute("href"))
      ) {
        node.setAttribute("xlink:show", "new");
      }
    });
    return {
      md,
      codeStyle,
      markdownStyle,
      DOMPurify,
    };
  });

  const renderedContent = computed(() => {
    if (!deps.value) return mergedOptions.placeholder;

    // we met a valid #{issue_id} in which issue_id is an integer and >= 0
    // render a link to the issue
    const formatted = unref(markdown)
      .split(/(#\d+)\b/)
      .map((part) => {
        if (!part.startsWith("#")) {
          return part;
        }
        const id = parseInt(part.slice(1), 10);
        if (!Number.isNaN(id) && id > 0) {
          const project = unref(projectRef);
          if (project) {
            // Here we assume that the referenced issue and the current issue are always
            // in the same project
            // if project is specified
            const path = `projects/${extractProjectResourceName(
              project.name
            )}/issues/${id}`;
            const url = `${window.location.origin}/${path}`;
            return `[${t("common.issue")} #${id}](${url})`;
          } else {
            return `[${t("common.issue")} #${id}](${
              window.location.origin
            }/issue/${id})`;
          }
        }
        return part;
      })
      .join("");
    const { md, DOMPurify, codeStyle, markdownStyle } = deps.value;
    const rendered = md.render(formatted);
    const html = DOMPurify.sanitize(rendered);

    return [
      `<head>`,
      `<style>${codeStyle}</style>`,
      `<style>${markdownStyle}</style>`,
      `</head>`,
      `<body style="overflow: auto;" class="markdown-body">`,
      html,
      `</body>`,
    ].join("\n");
  });

  const adjustIframe = () => {
    if (!iframeRef.value) return;

    nextTick(() => {
      if (!iframeRef.value) return;
      const height =
        iframeRef.value.contentDocument?.documentElement.offsetHeight ?? 0;
      const normalizedHeight = minmax(
        height,
        mergedOptions.minHeight,
        mergedOptions.maxHeight
      );
      iframeRef.value.style.height = `${normalizedHeight + 2}px`;
    });
  };

  const ready = computed(() => {
    return !!deps.value;
  });
  return { adjustIframe, ready, renderedContent };
};
