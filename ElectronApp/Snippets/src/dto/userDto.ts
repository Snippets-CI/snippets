export interface UserDto {
  user_id: string;
  username: string;
}

export function isUserDto(obj: any): obj is UserDto {
  return typeof obj.user_id === "string" && typeof obj.username === "string";
}
