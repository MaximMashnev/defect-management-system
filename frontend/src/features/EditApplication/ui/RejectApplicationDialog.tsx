import {
    TextField,
    Button,
    DialogTitle, 
    DialogContent, 
    DialogActions, 
    Dialog,
    styled
} from '@mui/material';
import { FC, useEffect, useState} from 'react';
import { ApplicationEntity } from '../../../entities/Application';
import applicationsStore from '../../../entities/Application/model/store';
import { ApplicationStatus } from '../../../entities/Application/model/types';
import { observer } from 'mobx-react-lite';

interface RejectApplicationDialogProps {
    open: boolean;
    data: ApplicationEntity;
    onClose: () => void;
}

const RejectForm = styled('form')({
    display: "flex",
    flexDirection: "column",
    gap: "1rem",
    paddingTop: "6px"
});

const RejectApplicationDialog: FC<RejectApplicationDialogProps> = observer(({open, data, onClose}) => {

    const [reason, setReason] = useState("");
    const [isLoading, setIsLoading] = useState(false);

    useEffect(() => {
        if (!open) return;
        setReason("");
    }, [open, data])

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsLoading(true);
        const application: ApplicationEntity = {
            ...data,
            status: ApplicationStatus.rejected,
            reason: reason,
            dateResult: new Date()
        }

        await applicationsStore.editApplication(application);
        if (applicationsStore.errorEdit) {
            window.alert(applicationsStore.errorEdit);
        }
        else {
            window.alert("Статус заявки обновлен.");
            onClose();
        }

        setIsLoading(false);
    }

    return (
      <Dialog open={open} onClose={onClose} fullWidth maxWidth="sm">
        <DialogTitle>{`Отклонение заявки "${data?.id}"`}</DialogTitle>
        <DialogContent>
            <RejectForm id="RejectForm" onSubmit={handleSubmit}>
                <TextField
                    value={reason}
                    onChange={e => setReason(e.target.value)}
                    label="Причина отклонения заявки"
                    multiline
                    maxRows={4}
                    required
                    type="text"
                    disabled={isLoading}
                /> 
            </RejectForm>
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose} disabled={isLoading}>Отмена</Button>
          <Button type="submit" form="RejectForm" disabled={isLoading} loading={isLoading}>
            Сохранить
          </Button>
        </DialogActions>
      </Dialog>
    )
})

export default RejectApplicationDialog;