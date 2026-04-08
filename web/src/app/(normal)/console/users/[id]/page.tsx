import { notFound } from 'next/navigation';
import { UserDetail } from '@/components/features/console/user-detail';

interface UserDetailPageProps {
  params: Promise<{
    id: string;
  }>;
}

export default async function UserDetailPage({ params }: UserDetailPageProps) {
  const { id } = await params;
  const userId = Number(id);

  if (!Number.isInteger(userId) || userId <= 0) {
    notFound();
  }

  return <UserDetail userId={userId} />;
}
