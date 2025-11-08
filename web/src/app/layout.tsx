import type { Metadata } from "next";
import "./globals.css"
import { Toaster } from "sonner"

export const metadata: Metadata = {
  title: "Horo Project",
  description: "Your personalized horoscope app",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className="font-sans antialiased bg-zinc-50 dark:bg-black">
        <main>{children}</main>
        <Toaster />
      </body>
    </html>
  );
}
