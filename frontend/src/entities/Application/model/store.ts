import { makeAutoObservable, runInAction } from "mobx";
import { ApplicationCategory, ApplicationEntity, NewApplication } from "./types";
import { Meta} from "../../../shared/api/services/types";
import { GridSortItem } from "@mui/x-data-grid/models/gridSortModel";
import { GridFilterModel, GridPaginationModel } from "@mui/x-data-grid";
import ApplicationService from "../services/applications.service";
import { FileDTO } from "../services/dto";

class ApplicationsStore {

    applications: ApplicationEntity[] = [];
    pagination: GridPaginationModel = {page: 0, pageSize: 5};
    meta: Meta | null = null;
    sort: GridSortItem | null = null;
    filter: GridFilterModel | null = null;
    categories: ApplicationCategory[] = [];

    isLoadingApplicationsForTable = false;
    isLoadingApplications = false;
    isLoadingApplication = false;
    isLoadingCategories = false;
    isAddingApplication = false;
    isEditingApplication = false;
    isUploadingFile = false;
    isLoadingFiles = false;

    errorApplicationsForTable: string | null = null;
    errorApplications: string | null = null;
    errorApplication: string | null = null; 
    errorCategories: string | null = null; 
    errorAdd: string | null = null; 
    errorEdit: string | null = null; 
    errorUpload: string | null = null;
    errorFiles: string | null = null;

    constructor () {
        makeAutoObservable(this, {}, {autoBind: true});
    }

    private handleError(error: unknown, defaultMessage: string): string {
        const message = error instanceof Error ? error.message : defaultMessage;
        console.error(message, error);
        return message;
    }

    async getApplicationsForTable() {
        this.isLoadingApplicationsForTable = true;
        this.errorApplicationsForTable = null;
        try {
            const data = await ApplicationService.getApplicationsForTable({
                sort: this.sort || undefined,
                filter: this.filter || undefined,
                pagination: this.pagination || {page: 0, pageSize: 5},
            });
            runInAction(() => {
                this.applications = data.items;
                this.meta = data.meta;
            });
        } 
        catch (error) {
            runInAction(() => {
                this.errorApplicationsForTable = this.handleError(error, "Ошибка загрузки заявок");
            });
        } finally {
            runInAction(() => {
                this.isLoadingApplicationsForTable = false;
            });
        }
    }

    async getApplications() {
        this.isLoadingApplications = true;
        this.errorApplications = null;
        try {
            const data = await ApplicationService.getApplications();
            runInAction(() => {
                this.applications = data;
            });
        }
        catch (error) {
            runInAction(() => {
                this.errorApplications = this.handleError(error, "Ошибка загрузки заявок");
            });
        } finally {
            runInAction(() => {
                this.isLoadingApplications = false;
            });
        }
    }

    async getApplication(id: number): Promise<ApplicationEntity | undefined> {
        this.isLoadingApplication = true;
        this.errorApplication = null;
        try {
            const data = await ApplicationService.getApplication(id);
            return data;
        }
        catch (error) {
            runInAction(() => {
                this.errorApplication = this.handleError(error, "Ошибка загрузки заявки");
            });
            return undefined;
        } finally {
            runInAction(() => {
                this.isLoadingApplication = false;
            });
        }
    }

    async getApplicationsCategories() {
        this.isLoadingCategories = true;
        this.errorCategories = null;
        try {
            const data = await ApplicationService.getApplicationsCategories();
            runInAction(() => {
                this.categories = data;
            });
        }
        catch (error) {
            runInAction(() => {
                this.errorCategories = this.handleError(error, "Ошибка при получении категорий");
            });
        } finally {
            runInAction(() => {
                this.isLoadingCategories = false;
            });
        }
    }

    async addApplication(newApplication: NewApplication): Promise<ApplicationEntity | undefined> {
        this.isAddingApplication = true;
        this.errorAdd = null;
        try {
            const data = await ApplicationService.addApplication(newApplication);
            return data;
        }
        catch (error) {
            runInAction(() => {
                this.errorAdd = this.handleError(error, "Ошибка создания заявки");
            });
            return undefined;
        } finally {
            runInAction(() => {
                this.isAddingApplication = false;
            });
        }
    }

    async editApplication(application: ApplicationEntity) {
        this.isEditingApplication = true;
        this.errorEdit = null;
        try {
            await ApplicationService.editApplication(application);
            await this.getApplicationsForTable();
        }
        catch (error) {
            runInAction(() => {
                this.errorEdit = this.handleError(error, "Ошибка редактирования заявки");
            });
        } finally {
            runInAction(() => {
                this.isEditingApplication = false;
            });
        }
    }

    async uploadFile(file: File): Promise<FileDTO | undefined> {
        this.isUploadingFile = true;
        this.errorUpload = null;
        try {
            const data = await ApplicationService.uploadFile(file);
            return data;
        }
        catch (error) {
            runInAction(() => {
                this.errorUpload = this.handleError(error, "Ошибка выгрузки файла");
            });
            return undefined;
        } finally {
            runInAction(() => {
                this.isUploadingFile = false;
            });
        }
    }

    async getFiles(idFiles: number[]): Promise<FileDTO[] | undefined> {
        this.isLoadingFiles = true;
        this.errorFiles = null;
        try {
            const data = await ApplicationService.getFiles(idFiles);
            return data;
        }
        catch (error) {
            runInAction(() => {
                this.errorFiles = this.handleError(error, "Ошибка загрузки файлов");
            });
            return undefined;
        } finally {
            runInAction(() => {
                this.isLoadingFiles = false;
            });
        }
    }

    setFilter(filter: GridFilterModel | null) {
        this.filter = filter;
    }

    setPagination(pagination: GridPaginationModel) {
        this.pagination = pagination;
    }

    setSort(sort: GridSortItem | null) {
        this.sort = sort;
    }

    resetErrors() {
        this.errorApplicationsForTable = null;
        this.errorApplications = null;
        this.errorApplication = null; 
        this.errorCategories = null; 
        this.errorAdd = null; 
        this.errorEdit = null; 
        this.errorUpload = null;
        this.errorFiles = null;
    }
}

const applicationsStore = new ApplicationsStore();
export default applicationsStore;