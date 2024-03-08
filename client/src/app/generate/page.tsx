"use client";
import QRCode from "react-qr-code";

const GenerateLink = () => {
  return (
    <div
      style={{ height: "auto", margin: "0 auto", maxWidth: 400, width: "100%" }}
    >
      <QRCode
        value="https://c9c5-2409-40f2-1041-f519-f5be-b512-7b31-4604.ngrok-free.app/generate/api"
        size={500}
        style={{ height: "auto", maxWidth: "100%", width: "100%" }}
        viewBox={`0 0 256 256`}
      />
    </div>
  );
};

export default GenerateLink;
