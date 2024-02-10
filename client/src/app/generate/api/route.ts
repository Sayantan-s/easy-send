import { openPDF } from "@/services/Scrapper";
import { NextApiHandler } from "next";

export const GET: NextApiHandler = (req, res) => {
  const pdfPath = "/public/resume_1YOE.pdf";
  openPDF(pdfPath);
  return new Response(JSON.stringify({ name: "Sayantan" }), {
    status: 200,
    statusText: "Yo nigga",
  });
};
