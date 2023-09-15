import { TextDocument } from "vscode-languageserver-textdocument";
import { CompletionParams } from "vscode-languageserver/browser";
import { LanguageState } from "../../types";
import { complete } from "./plsql";

export const completeV2 = async (
  params: CompletionParams,
  document: TextDocument,
  state: LanguageState
) => {
  const inputString = document.getText();
  const caretIndex = document.offsetAt(params.position);
  console.log("inputString", inputString);
  console.log("position", params.position);
  console.log("caretIndex", caretIndex);

  const result = await complete(inputString, caretIndex);

  return result;
};
