import { Box, styled } from "@mui/material";
import bgImage from "../../../shared/assets/imgs/bgImage.jpg";

const FormBox = styled(Box)({
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    justifyContent: "center",
    height: "100vh",
    position: "relative", 
    '&::before': {
        content: '""',
        position: "absolute",
        top: 0,
        left: 0,
        width: "100%",
        height: "100%",
        backgroundImage: `url(${bgImage})`,
        backgroundSize: "cover",
        backgroundPosition: "center",
        filter: "blur(2px)",
        zIndex: -1,
    }    
})

export default FormBox;