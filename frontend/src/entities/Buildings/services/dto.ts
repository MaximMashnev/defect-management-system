import { Meta } from "../../../shared/api/services/types";
import { BuildingEntity } from "../model/types";

export interface BuildingsDTO {
    meta: Meta;
    items: BuildingEntity[];
}