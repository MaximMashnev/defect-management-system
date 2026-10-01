import { httpService } from "../../../shared/api/services/http.service"
import config from "../../../shared/configs/config.json"
import { BuildingEntity, IncrementNumAppsBuilding, NewBuilding } from "../model/types";
import { BuildingsDTO } from "./dto";
import dayjs from "dayjs";
import { GridFilterModel, GridPaginationModel } from "@mui/x-data-grid";
import { GridSortItem } from "@mui/x-data-grid/models/gridSortModel";

interface TableParams {
    filter?: GridFilterModel;
    sort?: GridSortItem;
    pagination: GridPaginationModel;
}

const Building = config.api.endPoints.buildings;

const setDate = (building: BuildingEntity | NewBuilding) => {
    return {
        ...building,
        dateRegistration: dayjs(building.dateRegistration.setUTCHours(0, 0, 0, 0)).toISOString()
    }
};

const BuildingService = {
    getBuildingsForTable: async ({filter, sort, pagination}: TableParams) => {
        const params = new URLSearchParams();
        const filterItems = filter?.items[0];
        if (filterItems?.field && filterItems?.value) {
            const filterOperator = filterItems.operator === "is" ? "" : filterItems.operator;
            const filterValue = typeof filterItems.value === "object"
                ? dayjs(filterItems.value).toISOString()
                : String(filterItems.value);
            params.append(filterItems.field, `${filterOperator}${filterValue}`)
        }

        if (sort?.field && sort?.sort) {
            const sortDir = sort.sort === "desc" ? "-" : "";
            params.append("sortBy", `${sortDir}${sort.field}`);
        }
        
        params.append("page", `${pagination.page + 1}`);
        params.append("limit", `${pagination.pageSize}`);

        const { data } = await httpService.get<BuildingsDTO>(`${Building}?${params}`);
        return data;
    },

    getBuildings: async () => (await httpService.get<BuildingEntity[]>(Building)).data,

    addBuilding: async (newBuilding: NewBuilding) => 
        (await httpService.post<BuildingEntity>(Building, setDate(newBuilding))).data,

    editBuilding: async (editBuilding: BuildingEntity) => 
        (await httpService.patch<BuildingEntity>(`${Building}/${editBuilding.id}`, setDate(editBuilding))).data,

    incrementNumApps: async (building: IncrementNumAppsBuilding) => 
        (await httpService.patch<BuildingEntity>(`${Building}/${building.id}`, building)).data,

    deleteBuilding: async (deleteBuilding: BuildingEntity) =>
        (await httpService.delete<BuildingEntity>(`${Building}/${deleteBuilding.id}`)).data,

}

export default BuildingService