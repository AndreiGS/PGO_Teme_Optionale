import FileUpload from "@/components/ui/fileupload";

const HomePage = () => {
  return (
    <div className="min-h-screen bg-white grid place-items-center mx-auto py-8">
      <div className="text-blue-900 text-2xl font-bold flex flex-col items-center space-y-4">
        <FileUpload />
      </div>
    </div>
  )
}

export default HomePage;
