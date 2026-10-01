import { httpService } from "../../../shared/api/services/http.service"
import config from "../../../shared/configs/config.json"
import dayjs from "dayjs";
import { ApplicationCategory, ApplicationEntity, NewApplication } from "../model/types";
import { ApplicationsDTO, FileDTO } from "./dto";
import { GridFilterModel, GridPaginationModel } from "@mui/x-data-grid";
import { GridSortItem } from "@mui/x-data-grid/models/gridSortModel";

interface TableParams {
    filter?: GridFilterModel;
    sort?: GridSortItem;
    pagination: GridPaginationModel;
}

const Applications = config.api.endPoints.applications;
const ApplicationsCategories = config.api.endPoints.category;
const ApplicationsUploads = config.api.endPoints.uploads;

const ApplicationService = {
    getApplicationsForTable: async ({filter, sort, pagination}: TableParams) => {
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

        const { data } = await httpService.get<ApplicationsDTO>(`${Applications}?${params}`);
        return data;
    },

    getApplications: async () => (await httpService.get<ApplicationEntity[]>(`${Applications}`)).data,

    getApplication: async (id: number) => (await httpService.get<ApplicationEntity>(`${Applications}/${id}`)).data,

    getApplicationsCategories: async () => (await httpService.get<ApplicationCategory[]>(ApplicationsCategories)).data,

    addApplication: async (newApplication: NewApplication) => (await httpService.post<ApplicationEntity>(Applications, newApplication)).data,

    editApplication: async (editApplication: ApplicationEntity) => 
        (await httpService.patch<ApplicationEntity>(`${Applications}/${editApplication.id}`, editApplication)).data,

    uploadFile: async (file: File) => {
        const formData = new FormData();
        formData.append('file', file);
        const { data } = await httpService.post<FileDTO>(`${ApplicationsUploads}`, formData);
        return data;
    },

    getFiles: async (ids: number[]) => {
        const params = new URLSearchParams();
        ids.forEach(id => {
            params.append("id", `${id}`);
        })
        const { data } = await httpService.get<FileDTO[]>(`${ApplicationsUploads}?${params}`);
        return data;
    }
}

export default ApplicationService