import { makeAutoObservable, runInAction } from "mobx";
import { STORAGE_KEYS } from "../../../shared/consts";
import authService from "../services/auth.service";
import { AuthData } from "./types";
import { AuthUser } from "../../../entities/User/model/types";

class AuthStore {
    token: string | null = localStorage.getItem(STORAGE_KEYS.TOKEN) ?? null;
    user: AuthUser | null = JSON.parse(localStorage.getItem(STORAGE_KEYS.USER) || "null");

    isLoadingLogin = false;

    errorLogin: string | null = null;

    constructor() {
        makeAutoObservable(this, {}, {autoBind: true})
    }

    private handleError = (error: unknown, defaultMessage: string): string => {
        const message = error instanceof Error ? error.message : defaultMessage;
        console.error(message, error);
        return message;
    }

    setToken(token: string | null) {
        this.token = token;
        if (token) {
            localStorage.setItem(STORAGE_KEYS.TOKEN, token);
        } else {
            localStorage.removeItem(STORAGE_KEYS.TOKEN);
        }
    }

    setUser(user: AuthUser | null) {
        this.user = user;
        if (user) {
            localStorage.setItem(STORAGE_KEYS.USER, JSON.stringify(user));
        } else {
            localStorage.removeItem(STORAGE_KEYS.USER);
        }
    }

    async login(authData: AuthData): Promise<AuthUser | null> {
        this.isLoadingLogin = true;
        this.errorLogin = null
        try {
            const data = await authService.login(authData);
            runInAction(() =>{
                this.setUser(data.data);
                this.setToken(data.token)
            })
            return data.data;
        }
        catch (error) {
            runInAction(() => {
                this.errorLogin = this.handleError(error, "Ошибка авторизации");
            })
            return null;
        }
        finally {
            runInAction(() => {
                this.isLoadingLogin = false;
            })
        }
    }

    logout() {
        runInAction(() => {
            this.setToken(null);
            this.setUser(null);
        });
        localStorage.removeItem(STORAGE_KEYS.TOKEN);
        localStorage.removeItem(STORAGE_KEYS.USER);
    }
}

const authStore = new AuthStore();
export default authStore;