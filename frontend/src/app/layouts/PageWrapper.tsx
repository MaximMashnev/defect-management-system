import { Suspense, FC, useEffect } from "react";
import { Box, CircularProgress, styled } from "@mui/material";

interface PageWrapperProps {
    title: string;
    component: React.ComponentType;
}

const LoadingBox = styled(Box)({
    display: 'flex', 
    alignItems: "center", 
    justifyContent: "center", 
    height: "calc(100vh - 100px)"    
})

export const PageWrapper: FC<PageWrapperProps> = ({title, component: Component}) => {
    useEffect(() => {document.title = title}, [title]);

    return (
        <Suspense 
            fallback={
                <LoadingBox>
                    <CircularProgress size="4rem" />
                </LoadingBox>
            }>
            <Component />
        </Suspense>
    )
}