"use client"
import Card from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";

export default function HomePage() {
  const router = useRouter();
  return (
    <div className="relative min-h-screen">

      <div className="absolute top-4 right-4 p-4 flex gap-4">
        <Button 
          className="bg-blue-500 text-white hover:bg-blue-600"
          onClick={()=>{
            router.push('/signin')
          }}
        >sign in</Button>
        <Button 
          className="bg-blue-500 text-white hover:bg-blue-600"
          onClick={()=>{
            router.push('/signup')
          }}
        >sign up</Button>
      </div>

      <div className="flex min-h-screen items-center justify-center">
        <Card>
          <h1 className="text-5xl font-bold text-black dark:text-zinc-50 mb-4">
            Welcome to Horo
          </h1>
          <p className="text-lg text-zinc-600 dark:text-zinc-400 text-center max-w-md mx-auto">
            Your personalized horoscope app. Get started by signing in or creating an account.
          </p>
        </Card>
      </div>
    </div>
  );
}
