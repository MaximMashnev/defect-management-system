import { BuildingEntity } from "../../Buildings";

// export type ApplicationStatus = "новый" | "в работе" | "выполнено" | "отклонено";
export enum ApplicationStatus {
    new = "новый",
    inProgress = "в работе" ,
    completed = "выполнено",
    rejected = "отклонено"
}

export type BuildingID = Pick<BuildingEntity, "id">

export interface ApplicationEntity {
    id: number;
    name: string;
    description: string;
    email: string;
    dateSubmission: Date;
    dateInProgress?: Date;
    dateResult?: Date;
    status: ApplicationStatus;
    priority: number;
    building_id: number;
    upload_id: number[];
    user_id?: number;
    reason?: string;
    category: number;
}

export type NewApplication = Omit<ApplicationEntity, "id">

export interface ApplicationCategory {
    id: number;
    name: string;
    score: number;
}

export interface LevelCriticality {
    level: string;
    score: number;
}