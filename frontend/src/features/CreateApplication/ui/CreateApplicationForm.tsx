import { 
    Alert, 
    Button, 
    Container, 
    FormControl, 
    IconButton, 
    InputLabel, 
    MenuItem, 
    OutlinedInput, 
    Select, 
    styled, 
    TextField, 
    Typography, 
} from "@mui/material";
import { useEffect, useState } from "react";
import HomeIcon from '@mui/icons-material/Home';
import CloudUploadIcon from '@mui/icons-material/CloudUpload';
import { ApplicationCategory, ApplicationStatus, LevelCriticality, NewApplication } from "../../../entities/Application/model/types";
import applicationsStore from "../../../entities/Application/model/store";
import { LvlCriticality } from "../../../entities/Application/consts";
import { MAX_FILES, MAX_SIZE_MB, PATHS } from "../../../shared/consts";
import buildingsStore from "../../../entities/Buildings/model/store";
import { observer } from "mobx-react-lite";
import { BuildingEntity } from "../../../entities/Buildings";
import { FileDTO } from "../../../entities/Application/services/dto";
import { useNavigate } from "react-router-dom";
import { IncrementNumAppsBuilding } from "../../../entities/Buildings/model/types";
import { linkGeneration } from "../../../entities/Application/lib/linkHelpers";
import FormHeader from "../../../shared/ui/Form/FormHeader.styles";
import Form from "../../../shared/ui/Form/FormContainer.Styles";

interface EmailData {
    value: string;
    touched: boolean;
}

interface FormData {
    building: BuildingEntity | undefined;
    name: string;
    description: string;
    email: EmailData;
    category: ApplicationCategory | undefined;
    levelCriticality: LevelCriticality | undefined;
    files: FileList | undefined;
}

const VisuallyHiddenInput = styled('input')({
  clip: 'rect(0 0 0 0)',
  clipPath: 'inset(50%)',
  height: 1,
  overflow: 'hidden',
  position: 'absolute',
  bottom: 0,
  left: 0,
  whiteSpace: 'nowrap',
  width: 1,
});

const PageContainer = styled(Container)(({ theme }) => ({
    width: 600,
    backgroundColor: theme.palette.background.paper,
    borderRadius: "16px",
    padding: "12px"
}));

const CustomFormControl = styled(FormControl)({
    margin: 0
});

const CreateApplicationForm = observer(() => {

    const navigate = useNavigate();
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState("");
    const [formData, setFormData] = useState<FormData>(() => ({
        building: undefined,
        name: "",
        description: "",
        email: {value: "", touched: false},
        category: undefined,
        levelCriticality: undefined,
        files: undefined
    }));

    const isEmailValid = formData.email.value.split("@")[0].length > 4;
    const showEmailError = !isEmailValid && formData.email.touched;

    useEffect(() => {
        getData();
    }, [])

    const getData = async () => {
        await Promise.all([
            applicationsStore.getApplicationsCategories(),
            applicationsStore.getApplications(),
            buildingsStore.getBuildings(),
        ])
    }    

    const priorityСalculation = () => {
        let priority = formData.levelCriticality?.score ?? 0;
        priority += formData.category?.score ?? 0;

        if (formData.description.length >= 50) {
            priority += formData.description.length >= 200 ? 2 : 1;
        }

        if (formData.files?.length) {
            priority += formData.files.length >= 3 ? 2 : 1;
            priority += Array.from(formData.files).some(f => f.type.startsWith('video/')) ? 1 : 0;
        }

        const createdAt = new Date();
        priority += (createdAt.getHours() < 22 && createdAt.getHours() > 6) ? 0 : 2;

        priority += applicationsStore.applications.some(a => 
            a.building_id === formData.building?.id && a.category === formData.category?.id
        ) ? 2 : 0;

        priority += buildingsStore.buildings.some(b => 
            formData.building && b.numberApplications > formData.building.numberApplications
        ) ? 0 : 1;

        priority = Math.min(priority, 10);

        return {priority, createdAt};
    }

    const uploadedFiles = async (files: FileList) => {
        const uploadedFiles = await Promise.all(
            Array.from(files).map(async (value) => {
                try {
                    const fileData = await applicationsStore.uploadFile(value);
                    return fileData;
                }
                catch (error) {
                    console.log(error instanceof Error ? error.message : `Ошибка загрузки файла ${value.name}`);
                    return null;
                }
            }
        ));
        return uploadedFiles.filter((dto): dto is FileDTO => dto !== undefined && dto !== null).map(dto => dto.id);
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!formData.building || !formData.category || !formData.levelCriticality) {
            setError("Заполните все обязательные поля.");
            return;
        }

        setIsLoading(true);
        try {
            let filesData: number[] = [];
            if (formData.files && formData.files.length > 0) {
                filesData = await uploadedFiles(formData.files);
                if (filesData.length === 0) {
                    setIsLoading(false);
                    setError("Не удалось загрузить файлы, попробуйте ещё раз.")
                    return;
                } 
            }

            const priorityResult = priorityСalculation();

            const newApplication: NewApplication = {
                ...formData,
                email: formData.email.value,
                dateSubmission: priorityResult.createdAt,
                status: ApplicationStatus.new,
                building_id: formData.building.id,
                upload_id: filesData,
                priority: priorityResult.priority,
                category: formData.category.id
            }


            const data = await applicationsStore.addApplication(newApplication);
            if (data) {
                const incrementNumApps: IncrementNumAppsBuilding = {
                    ...formData.building,
                    numberApplications: formData.building.numberApplications + 1
                }                
                await buildingsStore.incrementNumApps(incrementNumApps);
                const link = linkGeneration(data);
                navigate(`../${PATHS.APPLICATION_TRACKING}/${link}`)                
            }
            else {
                throw new Error("Ошибка при привязки заявки к объекту");
            }
        }
        catch (error) {
            setError(error instanceof Error ? error.message : "Ошибка при создании заявки")
        }
        finally{
            setIsLoading(false);
        }
    }

    const handleFileInput = (files: FileList | null) => {
        if (!files) return;

        if (files.length > MAX_FILES) {
            setError(`Выбрано более ${MAX_FILES} файлов`)
            return;
        }

        if (Array.from(files).some(f => f.size > MAX_SIZE_MB * 1024 * 1024)) {
            setError(`Максимальный размер одного файла ${MAX_SIZE_MB} МБ`)
            return;
        }

        setFormData({...formData, files: files});
        setError("");
    }

    return (
        <>
            <PageContainer>
                <FormHeader>
                    <Typography variant="h6" component="span">
                        Создание заявки
                    </Typography>
                    <IconButton href={PATHS.MAIN} title="На главную">
                        <HomeIcon />
                    </IconButton>
                </FormHeader>
                {error && <Alert severity="error" sx={{marginBottom:"16px"}}>{error}</Alert>}
                <Form 
                    onSubmit={handleSubmit}
                >
                    <CustomFormControl>
                        <InputLabel>Строительный объект</InputLabel>
                        <Select
                            required
                            value={formData.building?.name}
                            onChange={e => setFormData({
                                ...formData, 
                                building: buildingsStore.buildings.find(b => b.id === Number(e.target.value))
                            })}
                            input={<OutlinedInput label="Строительный объект" />}
                            >
                            {buildingsStore.buildings ?
                                buildingsStore.buildings.map(b => (
                                    <MenuItem
                                        key={b.id}
                                        value={b.id}
                                    >
                                        {b.name}
                                    </MenuItem>
                                ))
                                :
                                <MenuItem
                                    disabled={true}
                                >
                                    Загрузка...
                                </MenuItem>
                            }
                        </Select>
                    </CustomFormControl>
                    <TextField
                        value={formData.name}
                        onChange={e => setFormData({...formData, name: e.target.value})}
                        label="Название проблемы"
                        required
                        type="text"
                        disabled={isLoading}
                    />       
                    <TextField
                        value={formData.description}
                        onChange={e => setFormData({...formData, description: e.target.value})}
                        label="Подробное описание проблемы"
                        multiline
                        maxRows={4}
                        required
                        type="text"
                        disabled={isLoading}
                    />         
                    <TextField
                        value={formData.email.value}
                        onChange={e => setFormData({
                            ...formData, 
                            email: {value: e.target.value, touched: true}
                        })}
                        error={showEmailError}
                        helperText={showEmailError && "Минимальная длина имени почты (до @) 4 символа"}
                        label="Email для обратной связи"
                        placeholder="example@ex.com"
                        required
                        type="email"
                        disabled={isLoading}
                    />
                    <CustomFormControl>
                        <InputLabel>Категория проблемы</InputLabel>
                        <Select
                            required
                            value={formData.category}
                            onChange={e => setFormData({
                                ...formData, 
                                category: applicationsStore.categories.find(c => c.id === Number(e.target.value))
                            })}
                            input={<OutlinedInput label="Категория проблемы" />}
                            >
                            {applicationsStore.categories ?
                                applicationsStore.categories.map((category) => (
                                    <MenuItem
                                        key={category.id}
                                        value={category.id}
                                    >
                                        {category.name}
                                    </MenuItem>
                                ))
                                :
                                <MenuItem
                                    disabled={true}
                                >
                                    Загрузка...
                                </MenuItem>
                            }

                        </Select>
                    </CustomFormControl>
                    <CustomFormControl>
                        <InputLabel>Уровень критичности</InputLabel>
                        <Select
                            required
                            value={formData.levelCriticality?.level}
                            onChange={e => setFormData({
                                ...formData, 
                                levelCriticality: LvlCriticality.find((item) => item.level === e.target.value)
                            })}
                            input={<OutlinedInput label="Уровень критичности" />}
                            >
                                {LvlCriticality.map((level) => (
                                    <MenuItem
                                        key={level.level}
                                        value={level.level}
                                    >
                                        {level.level}
                                    </MenuItem>
                                ))}
                        </Select>
                    </CustomFormControl>
                    <Button
                        component="label"
                        role={undefined}
                        variant="outlined"
                        tabIndex={-1}
                        startIcon={<CloudUploadIcon />}
                        >
                        Загрузить файлы
                        <VisuallyHiddenInput
                            type="file"
                            accept="image/*, video/*"
                            onChange={e => handleFileInput(e.target.files)}
                            multiple
                        />
                    </Button>
                    <Button 
                        type="submit" 
                        variant="contained" 
                        loading={isLoading}
                        disabled={isLoading || Boolean(error)}
                    >
                        Отправить заявку
                    </Button>                    
                </Form>    
            </PageContainer>    
        </>    
    )
}) 

export default CreateApplicationForm;