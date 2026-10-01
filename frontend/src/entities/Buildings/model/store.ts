import { makeAutoObservable, runInAction } from "mobx";
import { BuildingEntity, IncrementNumAppsBuilding, NewBuilding } from "./types";
import BuildingService from "../services/buildings.service";
import { Meta} from "../../../shared/api/services/types";
import { GridSortItem } from "@mui/x-data-grid/models/gridSortModel";
import { GridFilterModel, GridPaginationModel } from "@mui/x-data-grid";

class BuildingsStore {
    buildings: BuildingEntity[] = [];
    pagination: GridPaginationModel = {page: 0, pageSize: 5};
    meta: Meta | null = null;
    sort: GridSortItem | null = null;
    filter: GridFilterModel | null = null;

    isLoadingBuildingsForTable = false;
    isLoadingBuildings = false;
    isAddingBuilding = false;
    isEditingBuilding = false;
    isDeletingBuilding = false;
    isIncrementingNumApps = false;

    errorBuildingsForTable: string | null = null;
    errorBuildings: string | null = null;
    errorAdd: string | null = null;
    errorEdit: string | null = null;
    errorDelete: string | null = null;
    errorIncrementNumApps: string | null = null;

    constructor() {
        makeAutoObservable(this, {}, {autoBind: true});
    }

    private handleError(error: unknown, defaultMessage: string): string {
        const message = error instanceof Error ? error.message : defaultMessage;
        console.error(message, error);
        return message;
    }

    async getBuildingsForTable() {
        this.isLoadingBuildingsForTable = true;
        this.errorBuildingsForTable = null;
        try {
            const data = await BuildingService.getBuildingsForTable({
                sort: this.sort || undefined,
                pagination: this.pagination,
                filter: this.filter || undefined,
            });
            runInAction(() => {
                this.buildings = data.items;
                this.meta = data.meta;
            });
        } 
        catch (error) {
            runInAction(() => {
                this.errorBuildingsForTable = this.handleError(error, "Ошибка получения объектов");
            });
        }
        finally {
            runInAction(() => {
                this.isLoadingBuildingsForTable = false;
            });  
        }
    }

    async getBuildings() {
        this.isLoadingBuildings = true;
        this.errorBuildings = null;
        try {
            const data = await BuildingService.getBuildings();
            runInAction(() => {
                this.buildings = data;
            });
        } 
        catch (error) {
            runInAction(() => {
                this.errorBuildings = this.handleError(error, "Ошибка получения объектов");
            });
        }
        finally {
            runInAction(() => {
                this.isLoadingBuildings = false;
            });  
        }
    }

    async addBuilding(newObject: NewBuilding) {
        this.isAddingBuilding = true;
        this.errorAdd = null;
        try {
            await BuildingService.addBuilding(newObject);
            await this.getBuildingsForTable();
        }
        catch (error) {
            runInAction(() => {
                this.errorAdd = this.handleError(error, "Ошибка добавления объекта");
            });
        }
        finally {
            runInAction(() => {
                this.isAddingBuilding = false;
            });  
        }
    }

    async editBuilding(building: BuildingEntity) {
        this.isEditingBuilding = true;
        this.errorEdit = null;
        try {
            await BuildingService.editBuilding(building);
            await this.getBuildingsForTable();
        }
        catch (error) {
            runInAction(() => {
                this.errorEdit = this.handleError(error, "Ошибка редактирования объекта");
            });
        }
        finally {
            runInAction(() => {
                this.isEditingBuilding = false;
            });  
        }
    }

    async incrementNumApps(building: IncrementNumAppsBuilding) {
        this.isIncrementingNumApps = true;
        this.errorIncrementNumApps = null;
        try {
            await BuildingService.incrementNumApps(building);
        }
        catch (error) {
            runInAction(() => {
                this.errorIncrementNumApps = this.handleError(error, "Ошибка при увеличении количества заявок у объекта");
            });
        }
        finally {
            runInAction(() => {
                this.isIncrementingNumApps = false;
            });  
        }
    }

    async deleteBuilding(building: BuildingEntity) {
        this.isDeletingBuilding = true;
        this.errorDelete = null;
        try {
            await BuildingService.deleteBuilding(building);
            await this.getBuildingsForTable();
        }
        catch (error) {
            runInAction(() => {
                this.errorDelete = this.handleError(error, "Ошибка удаления объекта");
            });
        }
        finally {
            runInAction(() => {
                this.isDeletingBuilding = false;
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
        this.errorBuildingsForTable = null;
        this.errorBuildings = null;
        this.errorAdd = null;
        this.errorEdit = null;
        this.errorDelete = null;
        this.errorIncrementNumApps = null;
    }
}

const buildingsStore = new BuildingsStore();
export default buildingsStore;