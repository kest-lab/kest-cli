'use client';

import { useState } from 'react';
import { Link } from 'react-router-dom';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm, useWatch } from 'react-hook-form';
import * as z from 'zod';
import { Check, X, AlertCircle, Eye, EyeOff, Loader2 } from 'lucide-react';

import { cn } from '@/utils';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Checkbox } from '@/components/ui/checkbox';
import { Google, Apple, GitHub } from '@/components/ui/icons';
import { useAuthStore, authSelectors } from '@/store/auth-store';
import { useRegister } from '@/hooks/use-auth';
import { useT } from '@/i18n';

const registerSchema = z.object({
    name: z.string().min(1, 'nameRequired').min(2, 'nameMinLength').max(50, 'nameTooLong'),
    email: z.string().min(1, 'emailRequired').email('emailInvalid'),
    password: z.string().min(8, 'passwordMinLength').max(100, 'passwordTooLong')
        .regex(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]/, 'passwordInvalid'),
    confirmPassword: z.string().min(1, 'confirmPasswordRequired'),
    terms: z.boolean().refine(val => val === true, 'termsRequired'),
}).refine(data => data.password === data.confirmPassword, {
    message: 'passwordsDoNotMatch',
    path: ['confirmPassword'],
});

type RegisterFormData = z.infer<typeof registerSchema>;

// Password strength indicator
function PasswordStrengthIndicator({ password, t }: { password: string; t: any }) {
    const requirements = [
        { key: 'length', label: t('passwordReqLength'), check: password.length >= 8 },
        { key: 'case', label: t('passwordReqCase'), check: /[a-z]/.test(password) && /[A-Z]/.test(password) },
        { key: 'number', label: t('passwordReqNumber'), check: /\d/.test(password) },
        { key: 'special', label: t('passwordReqSpecial'), check: /[@$!%*?&]/.test(password) },
    ];

    const passedCount = requirements.filter(r => r.check).length;
    const strengthPercent = (passedCount / requirements.length) * 100;

    return (
        <div className="mt-2 p-2.5 rounded-lg bg-muted/40 border border-border/50 space-y-2">
            <div className="flex items-center gap-2">
                <div className="flex-1 h-1.5 rounded-full bg-muted overflow-hidden">
                    <div
                        className={cn(
                            "h-full rounded-full transition-all duration-500 ease-out",
                            passedCount === 0 && "bg-muted-foreground/30",
                            passedCount > 0 && passedCount < 3 && "bg-warning",
                            passedCount >= 3 && passedCount < 4 && "bg-info",
                            passedCount === 4 && "bg-success",
                        )}
                        style={{ width: `${strengthPercent}%` }}
                    />
                </div>
                <span className={cn(
                    "text-xs font-semibold tabular-nums",
                    passedCount === 4 ? "text-success" : "text-muted-foreground"
                )}>
                    {passedCount}/4
                </span>
            </div>

            <div className="grid grid-cols-2 gap-x-2 gap-y-1">
                {requirements.map((req) => (
                    <div
                        key={req.key}
                        className={cn(
                            "flex items-center gap-1.5 text-xs transition-colors",
                            req.check ? "text-success" : "text-muted-foreground"
                        )}
                    >
                        <div className={cn(
                            "flex items-center justify-center w-3.5 h-3.5 rounded-full transition-all",
                            req.check ? "bg-success/20" : "bg-muted"
                        )}>
                            {req.check ? <Check className="w-2 h-2" /> : <X className="w-2 h-2 opacity-50" />}
                        </div>
                        <span className="truncate">{req.label}</span>
                    </div>
                ))}
            </div>
        </div>
    );
}

export function RegisterForm({ className }: { className?: string }) {
    const t = useT('auth');
    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);

    const { mutateAsync: register, isPending: isMutationLoading } = useRegister();
    const systemFeatures = useAuthStore.use.systemFeatures();

    const canRegister = systemFeatures ? authSelectors.canRegister(useAuthStore.getState()) : false;
    const hasSocialLogin = systemFeatures ? authSelectors.hasSocialLogin(useAuthStore.getState()) : false;

    const {
        register: registerField,
        handleSubmit,
        formState: { errors, isSubmitting },
        control,
    } = useForm<RegisterFormData>({
        resolver: zodResolver(registerSchema),
        defaultValues: { name: '', email: '', password: '', confirmPassword: '', terms: false },
    });

    const password = useWatch({ control, name: 'password' }) || '';
    const isFormLoading = isMutationLoading || isSubmitting;

    const onSubmit = async (data: RegisterFormData) => {
        try { await register(data); } catch (err) { console.error(err); }
    };

    const getErrorMessage = (errorKey: string | undefined) => {
        if (!errorKey) return undefined;
        return t(errorKey) || errorKey;
    };

    const inputClass = "h-10 bg-muted/30 border-border/80 focus:border-primary focus:ring-2 focus:ring-primary/20 transition-all";

    if (systemFeatures && !canRegister) {
        return (
            <div className={cn('flex flex-col gap-4', className)}>
                <Card className="border-border/60 shadow-xl dark:shadow-2xl dark:shadow-black/20 bg-card/95 backdrop-blur-sm">
                    <CardHeader className="text-center px-6 pb-2 pt-5">
                        <CardTitle className="text-2xl font-bold">{t('registrationDisabled')}</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4 px-6 pb-5">
                        <Alert>
                            <AlertCircle className="h-4 w-4" />
                            <AlertDescription className="text-sm">{t('registrationDisabledMessage')}</AlertDescription>
                        </Alert>
                        <div className="text-center">
                            <Link to="/login" title={t('backToSignIn')} className="text-sm font-semibold text-primary hover:underline">
                                {t('backToSignIn')}
                            </Link>
                        </div>
                    </CardContent>
                </Card>
            </div>
        );
    }

    return (
        <div className={cn('flex flex-col gap-4', className)}>
            <Card className="border-border/60 shadow-xl dark:shadow-2xl dark:shadow-black/20 bg-card/95 backdrop-blur-sm">
                <CardHeader className="text-center space-y-0.5 px-6 pb-1 pt-5">
                    <CardTitle className="text-2xl font-bold">{t('createAccount')}</CardTitle>
                    <p className="text-sm text-muted-foreground">{t('getStarted')}</p>
                </CardHeader>

                <CardContent className="space-y-3 px-6 pb-5 pt-3">
                    <form onSubmit={handleSubmit(onSubmit)} className="space-y-2.5">
                        <div className="space-y-1">
                            <Label htmlFor="name" className="text-sm font-medium">{t('fullName')}</Label>
                            <Input id="name" type="text" placeholder={t('enterFullName')} autoComplete="name"
                                disabled={isFormLoading} className={inputClass} {...registerField('name')} />
                            {errors.name && <p className="text-xs text-destructive">{getErrorMessage(errors.name.message)}</p>}
                        </div>

                        <div className="space-y-1">
                            <Label htmlFor="email" className="text-sm font-medium">{t('email')}</Label>
                            <Input id="email" type="email" placeholder={t('enterEmail')} autoComplete="email"
                                disabled={isFormLoading} className={inputClass} {...registerField('email')} />
                            {errors.email && <p className="text-xs text-destructive">{getErrorMessage(errors.email.message)}</p>}
                        </div>

                        <div className="space-y-1">
                            <Label htmlFor="password" className="text-sm font-medium">{t('password')}</Label>
                            <div className="relative">
                                <Input id="password" type={showPassword ? 'text' : 'password'} placeholder={t('createPassword')}
                                    autoComplete="new-password" disabled={isFormLoading} className={cn(inputClass, "pr-10")} {...registerField('password')} />
                                <Button type="button" variant="ghost" size="sm" tabIndex={-1}
                                    className="absolute right-0 top-0 h-full px-3 hover:bg-transparent"
                                    onClick={() => setShowPassword(!showPassword)}>
                                    {showPassword ? <EyeOff className="h-4 w-4 text-muted-foreground" /> : <Eye className="h-4 w-4 text-muted-foreground" />}
                                </Button>
                            </div>
                            {password.length > 0 && <PasswordStrengthIndicator password={password} t={t} />}
                            {errors.password && password.length === 0 && <p className="text-xs text-destructive">{getErrorMessage(errors.password.message)}</p>}
                        </div>

                        <div className="space-y-1">
                            <Label htmlFor="confirmPassword" className="text-sm font-medium">{t('confirmPassword')}</Label>
                            <div className="relative">
                                <Input id="confirmPassword" type={showConfirmPassword ? 'text' : 'password'} placeholder={t('confirmYourPassword')}
                                    autoComplete="new-password" disabled={isFormLoading} className={cn(inputClass, "pr-10")} {...registerField('confirmPassword')} />
                                <Button type="button" variant="ghost" size="sm" tabIndex={-1}
                                    className="absolute right-0 top-0 h-full px-3 hover:bg-transparent"
                                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}>
                                    {showConfirmPassword ? <EyeOff className="h-4 w-4 text-muted-foreground" /> : <Eye className="h-4 w-4 text-muted-foreground" />}
                                </Button>
                            </div>
                            {errors.confirmPassword && <p className="text-xs text-destructive">{getErrorMessage(errors.confirmPassword.message)}</p>}
                        </div>

                        <div className="flex items-start gap-2.5 pt-0.5">
                            <Checkbox id="terms" disabled={isFormLoading} className="mt-0.5" {...registerField('terms')} />
                            <div className="space-y-0.5">
                                <Label htmlFor="terms" className="text-sm font-normal text-muted-foreground leading-relaxed">
                                    {t('agreeToTerms')}{' '}
                                    <Link to="/terms" title={t('termsOfService')} className="text-primary hover:underline" target="_blank">{t('termsOfService')}</Link>{' '}
                                    {t('and')}{' '}
                                    <Link to="/privacy" title={t('privacyPolicy')} className="text-primary hover:underline" target="_blank">{t('privacyPolicy')}</Link>
                                </Label>
                                {errors.terms && <p className="text-xs text-destructive">{getErrorMessage(errors.terms.message)}</p>}
                            </div>
                        </div>

                        <Button
                            type="submit"
                            className="w-full h-10 font-medium bg-linear-to-r from-primary to-primary-deeper hover:from-primary/90 hover:to-primary-deeper/90 shadow-md hover:shadow-lg transition-all duration-200 mt-1"
                            disabled={isFormLoading}
                        >
                            {isFormLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                            {t('signUp')}
                        </Button>
                    </form>

                    {hasSocialLogin && (
                        <>
                            <div className="relative">
                                <div className="absolute inset-0 flex items-center"><span className="w-full border-t border-border/50" /></div>
                                <div className="relative flex justify-center text-xs">
                                    <span className="bg-card px-3 text-muted-foreground">{t('orContinueWith')}</span>
                                </div>
                            </div>
                            <div className="grid grid-cols-3 gap-2">
                                <Button variant="outline" type="button" className="h-10 hover:bg-muted/50 hover:border-border transition-colors" disabled={isFormLoading}><Google className="h-4 w-4" /></Button>
                                <Button variant="outline" type="button" className="h-10 hover:bg-muted/50 hover:border-border transition-colors" disabled={isFormLoading}><Apple className="h-4 w-4" /></Button>
                                <Button variant="outline" type="button" className="h-10 hover:bg-muted/50 hover:border-border transition-colors" disabled={isFormLoading}><GitHub className="h-4 w-4" /></Button>
                            </div>
                        </>
                    )}

                    <div className="text-center text-sm text-muted-foreground">
                        {t('hasAccount')}{' '}
                        <Link to="/login" title={t('signIn')} className="font-semibold text-primary hover:underline transition-colors">{t('signIn')}</Link>
                    </div>
                </CardContent>
            </Card>
        </div>
    );
}
