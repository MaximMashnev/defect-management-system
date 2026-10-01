import { UserEntity } from "../../../entities/User";
import { AuthUser } from "../../../entities/User/model/types";

export interface LoginDataDTO {
    data: AuthUser;
    token: string;
}

export interface AuthData {
    login: string;
    password: string;
}

export type LoginForm = Pick<UserEntity, "login" | "password">;