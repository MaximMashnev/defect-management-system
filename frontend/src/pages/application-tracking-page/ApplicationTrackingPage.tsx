import { Link, useNavigate, useParams } from "react-router-dom"
import applicationsStore from "../../entities/Application/model/store";
import { 
    Accordion,
    AccordionDetails, 
    AccordionSummary, 
    Box, 
    Button, 
    Container, 
    IconButton, 
    ImageList, 
    ImageListItem, 
    styled, 
    Typography 
} from "@mui/material";
import { PATHS } from "../../shared/consts";
import HomeIcon from '@mui/icons-material/Home';
import { useEffect, useState } from "react";
import Image from "../../shared/assets/imgs/bgImage.jpg";
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import { ApplicationEntity } from "../../entities/Application";
import { FileDTO } from "../../entities/Application/services/dto";
import { ApplicationStatus } from "../../entities/Application/model/types";
import { linkDecrypting } from "../../entities/Application/lib/linkHelpers";
import FormHeader from "../../shared/ui/Form/FormHeader.styles";

const PageBox = styled(Box)({
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

const TrackingContainer = styled(Container)({
    boxSizing: "border-box",
    width: 600, 
    backgroundColor: "white", 
    borderRadius: "16px", 
    padding: "12px"
})

const ImageListBox = styled(ImageList)({
    width: 500, 
    maxHeight: 450 
})

export default function ApplicationsTrackingPage() {

    const navigate = useNavigate();
    const params = useParams<{id: string}>();
    const [app, setApp] = useState<ApplicationEntity | null>(null);
    const [imgs, setImgs] = useState<FileDTO[] | null>(null);
    const [isLoading, setIsloading] = useState(false);

    useEffect(() => {
        if (!params.id) {
            navigate(PATHS.MAIN, {replace: true});
        };
        getInfoApp();
    }, [params.id, navigate])

    const getInfoApp = async () => {
        setIsloading(true);
        if (!params.id) return null;
        const applicationId = linkDecrypting(params.id);
        if (!applicationId) return null;
        try {
           const app = await applicationsStore.getApplication(applicationId); 
            if (app) {
                setApp(app);
                if (app.upload_id.length > 0) {
                    const imgs = await applicationsStore.getFiles(app.upload_id);
                    if (imgs) setImgs(imgs);
                }
            }
        }
        catch (error) {
            const message = error instanceof Error ? error.message : "Ошибка загрузки заявки"
            console.error(message);
        }
        finally {
            setIsloading(false);
        }
    }

    return (
        <PageBox>
            <TrackingContainer>
                <FormHeader>
                    <Typography variant="h6" component="span">
                        {app ? `Отслеживание заявки ${params.id}` : "Заявка не найдена"}
                    </Typography>
                    <IconButton href={PATHS.MAIN} title="На главную">
                        <HomeIcon />
                    </IconButton>
                </FormHeader>
                {app ?
                    <>
                        <Typography variant="body1" component="span">
                            Текущей статус заявки: {app.status}
                            {app.reason && <><br />{`Причина: ${app.reason}`}</>}
                        </Typography>                    
                        <Accordion>
                            <AccordionSummary
                                expandIcon={<ExpandMoreIcon />}
                            >
                                <Typography component="span">
                                    История изменений статуса
                                </Typography>
                            </AccordionSummary>
                            <AccordionDetails>
                                {`Статус "${ApplicationStatus.new}": ${new Date(app.dateSubmission).toLocaleString()}`}
                                {app.dateInProgress && 
                                <><br />{`Статус "${ApplicationStatus.inProgress}": ${new Date(app.dateInProgress).toLocaleString()}`}</>}
                                {app.dateResult 
                                && <><br />{`Статус "${ApplicationStatus.completed}/${ApplicationStatus.rejected}": ${new Date(app.dateResult).toLocaleString()}`}</>}
                            </AccordionDetails>
                        </Accordion>
                        {
                            app.upload_id.length > 0 &&        
                            <Accordion>
                                <AccordionSummary
                                    expandIcon={<ExpandMoreIcon />}
                                >
                                    <Typography component="span">
                                        Прикрепленные файлы
                                    </Typography>
                                </AccordionSummary>
                                <AccordionDetails>
                                    {isLoading 
                                        ? "Загрузка..."
                                        : <ImageListBox cols={3} rowHeight={164}>
                                            {
                                                !imgs
                                                ? "Файлы не найдены"
                                                : imgs.map((item) => (
                                                    <ImageListItem key={item.id}>
                                                    <img
                                                        srcSet={`${item.url}?w=164&h=164&fit=crop&auto=format&dpr=2 2x`}
                                                        src={`${item.url}?w=164&h=164&fit=crop&auto=format`}
                                                        alt={item.fileName}
                                                        loading="lazy"
                                                    />
                                                    </ImageListItem>
                                                ))
                                            }
                                        </ImageListBox>                                    
                                    }
                                </AccordionDetails>
                            </Accordion>                   
                        }
                    </>
                    :
                    <Typography variant="body1" component="span">
                        {isLoading
                            ? "Загрузка..." 
                            : <Button component={Link} to={PATHS.MAIN}>На главную</Button>
                        }
                    </Typography>  
                }
            </TrackingContainer>
        </PageBox>
    )
}