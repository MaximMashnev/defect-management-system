export interface BuildingEntity {
    id: number;
    name: string;
    address: string;
    dateRegistration: Date;
    numberApplications: number;
}

export type NewBuilding = Omit<BuildingEntity, "id">;

export type IncrementNumAppsBuilding = Readonly<Omit<BuildingEntity, "numberApplications">> & Pick<BuildingEntity, "numberApplications">; 