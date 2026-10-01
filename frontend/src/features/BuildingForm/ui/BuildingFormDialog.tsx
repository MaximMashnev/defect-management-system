import {
    TextField,
    Button,
    DialogTitle, 
    DialogContent, 
    DialogActions, 
    Dialog,
} from '@mui/material';
import React, { FC, useEffect, useState} from 'react';
import { BuildingEntity } from '../../../entities/Buildings';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import dayjs from 'dayjs';
import { NewBuilding } from '../../../entities/Buildings/model/types';
import buildingsStore from '../../../entities/Buildings/model/store';
import { DATE_FORMAT } from '../../../shared/consts';
import Form from '../../../shared/ui/Form/FormContainer.Styles';
import { observer } from 'mobx-react-lite';

interface BuildingFormDialogProps {
    open: boolean;
    data: BuildingEntity | null;
    onClose: () => void;
}

interface FormData {
    name: string;
    address: string;
    dateRegistration: Date;
}

const BuildingFormDialog: FC<BuildingFormDialogProps> = observer(({open, data, onClose}) => {

    const [formData, setFormData] = useState<FormData>(() => ({
        name: data?.name ?? '',
        address: data?.address ?? '',
        dateRegistration: data?.dateRegistration ?? new Date(),
    }));
    const [isLoading, setIsLoading] = useState(false);

    useEffect(() => {
        if (!open) return;

        setFormData({
                name: data?.name ?? '',
                address: data?.address ?? '',
                dateRegistration: data?.dateRegistration ?? new Date(),
        });
    }, [open, data]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        
        if (!formData.dateRegistration) {
            window.alert("Выберите дату регистрации.");
            return;
        }

        setIsLoading(true);

        if (data) {
            const building: BuildingEntity = {...data, ...formData};
            await buildingsStore.editBuilding(building);
        }
        else {
            const building: NewBuilding = {...formData, numberApplications: 0};
            await buildingsStore.addBuilding(building);
        }
        if (buildingsStore.errorEdit || buildingsStore.errorAdd) {
            window.alert(buildingsStore.errorEdit || buildingsStore.errorAdd)
        }
        else {
            window.alert(data ? "Информация об объекте была обновлена." : "Новый объект был добавлен.");
            onClose();
        }

        setIsLoading(false);
    }

    return (
      <Dialog open={open} onClose={onClose}>
        <DialogTitle>{data ? `Редактирование объекта "${data.name}"` : "Создание объекта"}</DialogTitle>
        <DialogContent>
            <Form id="BuildingForm" onSubmit={handleSubmit}>
                <TextField
                    required
                    label="Название"
                    value={formData.name}
                    onChange={e => setFormData({...formData, name: e.target.value})}
                    type="text"
                    fullWidth
                    disabled={isLoading}
                />
                <TextField
                    required
                    label="Адрес"
                    value={formData.address}
                    onChange={e => setFormData({...formData, address: e.target.value})}
                    type="text"
                    fullWidth
                    disabled={isLoading}
                />
                <LocalizationProvider dateAdapter={AdapterDayjs}>
                    <DatePicker
                    label="Дата регистрации"
                    value={dayjs(formData.dateRegistration)}
                    onChange={newDate => setFormData({...formData, dateRegistration: dayjs(newDate).isValid() ? dayjs(newDate).toDate() : new Date()})}
                    format={DATE_FORMAT}
                    disabled={isLoading}
                    />
                </LocalizationProvider>
            </Form>
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose} disabled={isLoading}>Отмена</Button>
          <Button type="submit" form="BuildingForm" disabled={isLoading} loading={isLoading}>
            {data ? "Сохранить" : "Добавить"}
          </Button>
        </DialogActions>
      </Dialog>
    )
})

export default BuildingFormDialog;