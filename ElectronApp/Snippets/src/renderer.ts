import { SnippetDto } from "./dto/snippetDto";
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
import "bootstrap-select";
import "bootstrap-select/dist/css/bootstrap-select.min.css";
import axios from "axios";
import * as monaco from "monaco-editor";
import * as snippet from "./dto/snippetDto";
import * as user from "./dto/userDto";

import $ from "jquery";

/* ********************
 * Declarations
 * ********************/

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

/* ********************
 * Functions
 * ********************/

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

/**
 * @param {String} HTML representing a single element
 * @return {Element}
 */
function htmlToElement(html: string): ChildNode {
  const template = document.createElement("template");
  html = html.trim(); // Never return a text node of whitespace as the result
  template.innerHTML = html;
  return template.content.firstChild;
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
  });
}

function loadSnippets(): void {
  axios
    .get("http://localhost:8010/users/1/snippets")
    .then((response) => {
      // check and cast to Snippet
      if (!snippet.isSnippetDto(response.data[0])) {
        console.error("Invalid request");
        console.info(response);
      } else {
        const peopleArray: SnippetDto[] = Object.keys(response.data).map(
          (i) => response.data[i]
        );

        const ul = document.getElementById("snippetList");

        for (const data of peopleArray) {
          const html = `<li class="nav-item">
              <a class="nav-link" href="#">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
                  stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
                  class="feather feather-file-text"
                  style="--darkreader-inline-fill:none; --darkreader-inline-stroke:currentColor;"
                  data-darkreader-inline-fill="" data-darkreader-inline-stroke="">
                  <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                  <polyline points="14 2 14 8 20 8"></polyline>
                  <line x1="16" y1="13" x2="8" y2="13"></line>
                  <line x1="16" y1="17" x2="8" y2="17"></line>
                  <polyline points="10 9 9 9 8 9"></polyline>
                </svg>
                ${data.name}
              </a>
            </li>`;

          const li = htmlToElement(html);

          //const li = document.createElement("li");
          //li.appendChild(document.createTextNode(data.id));
          ul.appendChild(li);
        }
      }
    })
    .catch((error) => {
      console.log("Error while loading snippets: ", error);
    });
}

function loadSnippet(): void {
  axios
    .get("http://localhost:8010/users/1/snippets/3")
    .then((response) => {
      // check and cast to Snippet
      if (!snippet.isSnippetDto(response.data)) {
        console.error("Invalid request");
        console.info(response);
      } else {
        model.setValue(response.data.code);
        selector.selectedIndex = languages.indexOf(response.data.lang);
        monaco.editor.setModelLanguage(model, response.data.lang);
        document.getElementById("monacoSnippetName").textContent =
          response.data.name;
      }
    })
    .catch((error) => {
      console.log("Error while loading snippet: ", error);
    });
}

function loadUser(): void {
  axios
    .get("http://localhost:8010/users/1")
    .then((response) => {
      // check and cast to Snippet
      if (!user.isUserDto(response.data)) {
        console.error("Invalid request");
        console.info(response);
      } else {
        document.getElementById("userNameLink").textContent =
          "Hi " + response.data.name;
      }
    })
    .catch((error) => {
      console.log("Error while loading user: ", error);
    });
}

function loadMainApplication(): void {
  initializeMonacoEditor();
  loadUser();
  loadSnippets();
  loadSnippet();
}

function addLoginListener(): void {
  $("#loginBtn").click(() => {
    loadMainApplication();
    $("#loginModal").modal("hide").data("#loginModal", null);
  });

  $("#registerBtn").click(() => {
    loadMainApplication();
    $("#loginModal").modal("hide").data("#loginModal", null);
  });
}

// Load Login Modal on start up
$(window).on("load", function () {
  $("#loginModal").modal("show");
});

window.addEventListener("resize", updateDimensions.bind(this));
window.onclose = function (): void {
  window.removeEventListener("resize", updateDimensions.bind(this));
};

addLoginListener();
loadLanguages();
