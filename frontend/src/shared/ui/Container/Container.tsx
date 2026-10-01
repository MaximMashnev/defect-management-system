import { FC } from "react";
import { Container } from "./Container.styles";

const ContainerComponent: FC<{children: React.ReactNode}> = ({children}) => {
    return (
        <Container>
            {children}
        </Container>
    );
}

export default ContainerComponent