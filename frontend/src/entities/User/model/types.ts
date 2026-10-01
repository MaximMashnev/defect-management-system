export interface UserEntity {
    id: number;
    login: string;
    password: string;
};

export type AuthUser = Omit<UserEntity, "password">;