import Communication from "@/services/Communication";
import { NextApiHandler } from "next";

export const GET: NextApiHandler = (req, res) => {
  return Communication.response.success({
    status: 200,
    data: "Hello world",
  });
};
