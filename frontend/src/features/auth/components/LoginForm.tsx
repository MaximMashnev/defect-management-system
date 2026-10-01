import { Alert, Button, TextField, Typography, IconButton} from "@mui/material";
import Container from "../../../shared/ui/Container";
import { useState } from "react";
import { useNavigate } from "react-router";
import { PATHS } from "../../../shared/consts";
import HomeIcon from '@mui/icons-material/Home';
import FormBox from "../../../shared/ui/Form/FormBox.styles";
import FormHeader from "../../../shared/ui/Form/FormHeader.styles";
import Form from "../../../shared/ui/Form/FormContainer.Styles";
import authStore from "../model/store";
import { AuthData } from "../model/types";
import { observer } from "mobx-react-lite";

const LoginForm = observer(() => {

    const navigate = useNavigate();

    const [loginData, setLoginData] = useState<AuthData>(() => ({
        login: authStore.user?.login ?? '',
        password: ''
    }));
    const [touched, setTouched] = useState({login: false, password: false})

    const isLoginValid = loginData.login.length >= 4;
    const isPasswordValid = loginData.password.length >= 8;

    const showLoginError = touched.login && !isLoginValid;
    const showPasswordError = touched.password && !isPasswordValid;

    const authUser = async () => {
        if (!loginData) return;
        const data = await authStore.login(loginData);
        if (data) {
            navigate(`../${PATHS.TOOLPAD}`);
        }
    }   

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        authUser();
    }

    return (
        <FormBox>
            <Container>
                <FormHeader>
                    <Typography variant="h6" component="span">
                        Авторизация
                    </Typography>
                        <IconButton onClick={() => navigate(PATHS.MAIN)} title="На главную">
                            <HomeIcon />
                        </IconButton>
                </FormHeader>
                {authStore.errorLogin && <Alert severity="error">{authStore.errorLogin}</Alert>}
                <Form onSubmit={handleSubmit}>
                    <TextField
                        value={loginData.login}
                        onChange={e => setLoginData({...loginData, login: e.target.value})}
                        label="Логин"
                        placeholder="example@ex.com"
                        onBlur={() => {setTouched(prev => ({ ...prev, login: true}))}}
                        error={showLoginError}
                        helperText={showLoginError && "Минимальная длина 4 символа"}
                        type="email"
                        disabled={authStore.isLoadingLogin}
                    />
                    <TextField
                        value={loginData.password}
                        onChange={e => setLoginData({...loginData, password: e.target.value})}
                        label="Пароль"
                        placeholder="Password"
                        onBlur={() => {setTouched(prev => ({ ...prev, password: true}))}}
                        error={showPasswordError}
                        helperText={showPasswordError && "Минимальная длина 8 символов"}
                        type="password"
                        disabled={authStore.isLoadingLogin}
                    />
                    <Button 
                        type="submit" 
                        variant="contained" 
                        loading={authStore.isLoadingLogin}
                        disabled={!isLoginValid || !isPasswordValid}
                    >
                        Войти
                    </Button>                    
                </Form>    
            </Container>
        </FormBox>
    )
})

export default LoginForm;