/**
 * This file will automatically be loaded by webpack and run in the "renderer" context.
 * To learn more about the differences between the "main" and the "renderer" context in
 * Electron, visit:
 *
 * https://electronjs.org/docs/tutorial/application-architecture#main-and-renderer-processes
 *
 */

import "./index.css";
import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap";
import 'bootstrap-select'
import 'bootstrap-select/dist/css/bootstrap-select.min.css'
import axios from "axios";
import * as monaco from "monaco-editor";

const languages = [
  "abap",
  "apex",
  "azcli",
  "bat",
  "cameligo",
  "clojure",
  "coffee",
  "cpp",
  "csharp",
  "csp",
  "css",
  "dockerfile",
  "fsharp",
  "go",
  "graphql",
  "handlebars",
  "html",
  "ini",
  "java",
  "javascript",
  "json",
  "kotlin",
  "less",
  "lua",
  "markdown",
  "mips",
  "msdax",
  "mysql",
  "objective-c",
  "pascal",
  "pascaligo",
  "perl",
  "pgsql",
  "php",
  "postiats",
  "powerquery",
  "powershell",
  "pug",
  "python",
  "r",
  "razor",
  "redis",
  "redshift",
  "restructuredtext",
  "ruby",
  "rust",
  "sb",
  "scheme",
  "scss",
  "shell",
  "solidity",
  "sophia",
  "sql",
  "st",
  "swift",
  "tcl",
  "twig",
  "typescript",
  "vb",
  "xml",
  "yaml",
];

/*
axios.get("https://api.github.com/users/mapbox").then((response) => {
  console.log(response.data);
  console.log(response.status);
  console.log(response.statusText);
  console.log(response.headers);
  console.log(response.config);
});*/

const editor = monaco.editor.create(
  document.getElementById("monaco-container"),
  {
    theme: "vs-dark",
    scrollBeyondLastLine: false,
  }
);

const model = monaco.editor.createModel("function () {}", "javascript");
const selector = document.getElementById(
  "languageSelector"
) as HTMLSelectElement;

function updateDimensions(): void {
  editor.layout();
}

function initializeMonacoEditor(): void {
  editor.setModel(model);

  // Add bindings
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const myBinding = editor.addCommand(
    monaco.KeyMod.CtrlCmd | monaco.KeyCode.KEY_S,
    function () {
      alert("CTRL + S pressed! Save work");
    }
  );

  editor.layout();
}

function loadLanguages(): void {
  let count = 0;

  for (const language of languages) {
    const opt = document.createElement("option");
    opt.value = count.toString();
    opt.text = language;

    selector.appendChild(opt);
    count += 1;
  }

  selector.addEventListener("change", function () {
    const language = languages[parseInt(selector.value)];

    monaco.editor.setModelLanguage(model, language);
    console.log(language);
  });

  // manually set selected language for now
  selector.selectedIndex = languages.indexOf("javascript");
}

window.addEventListener("resize", updateDimensions.bind(this));
window.onclose = function (): void {
  window.removeEventListener("resize", updateDimensions.bind(this));
};

initializeMonacoEditor();
loadLanguages();
