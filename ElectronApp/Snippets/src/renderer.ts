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
// eslint-disable-next-line import/no-unresolved
import * as monaco from "monaco-editor";
import * as snippet from "./dto/snippetDto";
import * as user from "./dto/userDto";

import $ from "jquery";

/* ********************
 * Declarations
 * ********************/

const restApiConnectionString = "http://localhost:8010/";
const defaultLanguage = "markdown";
let currentUser: user.UserDto = null;

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

const model = monaco.editor.createModel("", "markdown");

const selector = $("#languageSelector").get(0) as HTMLSelectElement;

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

  // set default language
  selector.selectedIndex = languages.indexOf(defaultLanguage);
}

function setModelWithLanguage(loadedSnippet: snippet.SnippetDto): void {
  model.setValue(loadedSnippet.code);
  monaco.editor.setModelLanguage(model, loadedSnippet.language);
  selector.selectedIndex = languages.indexOf(loadedSnippet.language);
  $(".selectpicker").selectpicker("refresh");
  $("#monacoSnippetName").text(loadedSnippet.title);
}

async function loadUserAsync(
  connectionString: string,
  usermail: string
): Promise<user.UserDto> {
  return axios
    .post(`${connectionString}login`, {
      // eslint-disable-next-line @typescript-eslint/camelcase
      user_id: "",
      username: "",
      mail: usermail,
      password: "",
    })
    .then((response) => {
      if (!user.isUserDto(response.data)) {
        console.error(
          `Invalid request - expected UserDTO got:\n${response.data}`
        );
        console.info(response);
        return null;
      } else {
        return response.data;
      }
    })
    .catch((error) => {
      console.log(`Error while loading user ${usermail}`, error);
      return null;
    });
}

async function loadSnippetAsync(
  connectionString: string,
  userId: string,
  snippetId: string
): Promise<snippet.SnippetDto> {
  return axios
    .get(`${connectionString}user/${userId}/snippets/${snippetId}`)
    .then((response) => {
      if (!snippet.isSnippetDto(response.data)) {
        console.error(`Invalid request - expected SnippetDTO`);
        console.info(response);
        return null;
      } else {
        return response.data;
      }
    })
    .catch((error) => {
      console.log(
        `Error while loading snippet with id ${snippetId} for user ${userId}`,
        error
      );
      return null;
    });
}

async function loadSnippetsAsync(
  connectionString: string,
  userId: string
): Promise<snippet.SnippetDto[]> {
  return axios
    .get(`${connectionString}user/${userId}/snippets`)
    .then((response) => {
      if (response.data.length > 0) {
        if (!snippet.isSnippetDto(response.data[0])) {
          console.error(`Invalid request - expected SnippetDTO`);
          console.info(response);
          return [] as SnippetDto[];
        } else {
          const snippets: SnippetDto[] = Object.keys(response.data).map(
            (i) => response.data[i]
          );
          return snippets;
        }
      }
    })
    .catch((error) => {
      console.log(`Error while loading snippets for user ${userId}`, error);
      return [] as SnippetDto[];
    });
}

function createSnippetLinks(
  connectionString: string,
  snippets: SnippetDto[]
): void {
  const ul = document.getElementById("snippetList");

  for (const s of snippets) {
    const html = `<li class="nav-item">
                    <a class="nav-link" href="#">
                      ${s.title}
                    </a>
                  </li>`;

    const li = htmlToElement(html);
    li.addEventListener("click", async () => {
      loadSnippetAsync(
        connectionString,
        currentUser.user_id,
        s.snippet_id
      ).then((response) => {
        if (response != null) {
          setModelWithLanguage(response);
        }
      });
    });

    ul.appendChild(li);
  }
}

function loadMainApplication(usermail: string): void {
  loadUserAsync(restApiConnectionString, usermail).then((response) => {
    if (response != null) {
      currentUser = response;

      $("#loginModal").modal("hide").data("#loginModal", null);
      $("#userNameLink").text("Hi " + currentUser.username);

      loadSnippetsAsync(restApiConnectionString, currentUser.user_id).then(
        (response) => {
          if (response.length > 0) {
            createSnippetLinks(restApiConnectionString, response);
          }
        }
      );
    }
  });
}

function addLoginListener(): void {
  // Handle login
  $("#loginBtn").click(() => {
    const userName = $("#username").val() as string;
    const password = $("#password").val();

    loadMainApplication(userName);
  });

  // Handle Register
  $("#registerBtn").click(() => {
    let userName = $("#username").val();
    let password = $("#password").val();

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
initializeMonacoEditor();
loadLanguages();
