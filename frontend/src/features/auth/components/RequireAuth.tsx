import { FC, useEffect} from "react";
import { PATHS, STORAGE_KEYS } from "../../../shared/consts";
import { useNavigate } from "react-router-dom";
import authStore from "../model/store";
import { observer } from "mobx-react-lite";

interface RequireAuthProps {
    component: React.ComponentType;
}

export const RequireAuth: FC<RequireAuthProps> = observer(({component: Component}) => {

    const navigate = useNavigate();
    const token = authStore.token;

    useEffect(() => {
        if (!token) {
            navigate(PATHS.LOGIN);
        }
    }, [token, navigate])

    if (!token) {
        return null;
    }

    return <Component />;
})