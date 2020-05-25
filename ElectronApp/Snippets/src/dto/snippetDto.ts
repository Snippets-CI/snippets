export interface SnippetDto {
  id: string;
  name: string;
  lang: string;
  about: string;
  code: string;
}

export function isSnippetDto(obj: any): obj is SnippetDto {
  return (
    typeof obj.id === "string" &&
    typeof obj.name === "string" &&
    typeof obj.lang === "string" &&
    typeof obj.about === "string" &&
    typeof obj.code === "string"
  );
}
