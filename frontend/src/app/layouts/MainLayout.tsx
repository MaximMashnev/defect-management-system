import {Typography, Toolbar, Button, Box, AppBar, styled} from "@mui/material";
import { Link, Outlet } from "react-router-dom";
import LoginIcon from '@mui/icons-material/Login';
import LogoutIcon from '@mui/icons-material/Logout';
import AppsIcon from '@mui/icons-material/Apps';
import { PATHS, STORAGE_KEYS } from "../../shared/consts";
import authStore from "../../features/auth/model/store";
import { observer } from "mobx-react-lite";

const MainPageBox = styled(Box)({
  display: 'flex', 
  flexDirection: 'column', 
  height: '100vh'   
})

const MainLayout = observer(() => {
  return(
    <>
      <MainPageBox>
        <AppBar position="static">
          <Toolbar>
            <Typography variant="h6" component="div" sx={{ flexGrow: 1}}>
              <Link to={PATHS.MAIN} style={{textDecoration: "none", color: "white"}}>СтроительнаяКомпания</Link>
            </Typography>
            {
              authStore.user ? (
                <>
                  <Button color="inherit" endIcon={<AppsIcon />} component={Link} to={`../${PATHS.TOOLPAD}`}>
                    Панель инструментов
                  </Button >
                  <Button color="inherit" endIcon={<LogoutIcon />} onClick={authStore.logout}>
                    Выйти
                  </Button >
                </>
              ) : (
                <Button color="inherit" endIcon={<LoginIcon />} component={Link} to={PATHS.LOGIN}>
                  Авторизоваться
                </Button >
              )
            }
          </Toolbar>
        </AppBar>
        <Box component="main" sx={{ flex: 1}}>
          <Outlet />
        </Box>
      </MainPageBox>         
    </>
  )
})

export default MainLayout;