
import { AuthGuard } from '@/components/auth-guard';

export default function NormalLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <AuthGuard>{children}</AuthGuard>;
}
