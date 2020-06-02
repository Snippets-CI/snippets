export interface SnippetDto {
  snippet_id: string;
  title: string;
  language: string;
  category: string;
  code: string;
}

export function isSnippetDto(obj: any): obj is SnippetDto {
  return (
    typeof obj.snippet_id === "string" &&
    typeof obj.title === "string" &&
    typeof obj.language === "string" &&
    typeof obj.category === "string" &&
    typeof obj.code === "string"
  );
}
