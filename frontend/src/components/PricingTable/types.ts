import { Plan } from "@/types";

export interface LocalPlan extends Plan {
  label: string;
  image: string;
  buttonText: string;
  highlight: boolean;
  isFreePlan: boolean;
  isAvailable: boolean;
  pricing: string;
  priceUnit: string;
  priceInfo: string;
}
