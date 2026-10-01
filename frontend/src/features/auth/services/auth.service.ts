import config from "../../../shared/configs/config.json";
import { httpService } from "../../../shared/api/services/http.service";
import { LoginDataDTO, LoginForm } from "../model/types";
import { AuthUser } from "../../../entities/User/model/types";

const authMe = config.api.endPoints.auth_me;
const auth = config.api.endPoints.auth;

const authService = {
    login: async (loginForm: LoginForm) => (await httpService.post<LoginDataDTO>(auth, loginForm)).data,

    authMe: async () => (await httpService.get<AuthUser>(authMe)).data,
}

export default authService;