import { Box, Button, styled, Typography } from "@mui/material";
import { Link } from "react-router-dom";
import { PATHS } from "../../shared/consts";

const NotFoundBox = styled(Box)({
    display: "flex",
    flexDirection: "column",
    gap: "1rem",
    alignItems: "center",
    justifyContent: "center",
    height: "100vh"    
})

export default function NotFoundPage() {
    return (
        <NotFoundBox>
            <Typography variant="h5" component="span">
                Страница не найдена
            </Typography>
            <Button variant="contained" component={Link} to={PATHS.MAIN} sx={{ width: 300}}>На главную</Button>
        </NotFoundBox>
    )
}