import { Box, styled } from "@mui/material";
import CreateApplicationForm from "../../features/CreateApplication/ui/CreateApplicationForm";
import Image from "../../shared/assets/imgs/bgImage.jpg";

const CreatingPageBox = styled(Box)({
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    justifyContent: "center",
    height: "calc(100vh - 64px)",
    position: "relative", 
    '&::before': {
        content: '""',
        position: "absolute",
        top: 0,
        left: 0,
        width: "100%",
        height: "100%",
        backgroundImage: `url(${Image})`,
        backgroundSize: "cover",
        backgroundPosition: "center",
        filter: "blur(2px)",
        zIndex: -1,
    }    
})

export default function ApplicationCreatingPage() {
    return (
        <CreatingPageBox>
            <CreateApplicationForm />
        </CreatingPageBox>
    )
}