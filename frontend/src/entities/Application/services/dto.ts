import { Meta } from "../../../shared/api/services/types";
import { ApplicationEntity } from "../model/types";

export interface ApplicationsDTO {
    meta: Meta;
    items: ApplicationEntity[];
}

export interface FileDTO {
    bytes: number;
    fileName: string;
    format: string;
    height: number;
    id: number;
    url: string;
    width: number;
}