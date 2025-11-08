import { JwtPayload } from "jwt-decode";

export interface FirebaseClaims extends JwtPayload {
  user_id: string;
  email: string;
  name?: string;
  role?: string;
  picture?: string;
}
