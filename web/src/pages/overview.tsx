import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui';
import { useProviders, useRoutes, useProjects, useProxyRequests } from '@/hooks/queries';
import { Activity, Server, Route, FolderKanban, Zap, ArrowRight, CheckCircle, XCircle, Ban, LayoutDashboard } from 'lucide-react';
import { Link } from 'react-router-dom';

export function OverviewPage() {
  const { data: providers } = useProviders();
  const { data: routes } = useRoutes();
  const { data: projects } = useProjects();
  const { data: requests } = useProxyRequests({ limit: 10 });

  const stats = [
    { label: 'Providers', value: providers?.length ?? 0, icon: Server, color: 'text-info', href: '/providers' },
    { label: 'Routes', value: routes?.length ?? 0, icon: Route, color: 'text-accent', href: '/routes/claude' },
    { label: 'Projects', value: projects?.length ?? 0, icon: FolderKanban, color: 'text-warning', href: '/projects' },
    { label: 'Recent Requests', value: requests?.length ?? 0, icon: Activity, color: 'text-success', href: '/requests' },
  ];

  const completedRequests = requests?.filter((r) => r.status === 'COMPLETED').length ?? 0;
  const failedRequests = requests?.filter((r) => r.status === 'FAILED').length ?? 0;
  const cancelledRequests = requests?.filter((r) => r.status === 'CANCELLED').length ?? 0;
  const hasProviders = (providers?.length ?? 0) > 0;

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="h-[73px] flex items-center justify-between p-lg border-b border-border bg-surface-primary flex-shrink-0">
        <div className="flex items-center gap-md">
          <LayoutDashboard size={24} className="text-text-primary" />
          <div>
            <h2 className="text-headline font-semibold text-text-primary">Dashboard</h2>
            <p className="text-caption text-text-muted">Overview of your proxy gateway</p>
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto p-lg">
        <div className="space-y-xl animate-fade-in">
      {/* Welcome Section */}
      {!hasProviders && (
        <div className="text-center py-xxl">
          <div className="w-20 h-20 rounded-2xl bg-accent/10 flex items-center justify-center mx-auto mb-lg">
            <Zap size={40} className="text-accent" />
          </div>
          <h1 className="text-title1 font-bold text-text-primary mb-md">Welcome to Maxx Next</h1>
          <p className="text-body text-text-secondary max-w-md mx-auto mb-xl">
            AI API Proxy Gateway - Route your AI requests through multiple providers with intelligent failover and load balancing.
          </p>
          <Link
            to="/providers"
            className="inline-flex items-center gap-2 bg-accent text-white px-xl py-md rounded-lg hover:bg-accent-hover transition-colors"
          >
            Get Started
            <ArrowRight className="h-4 w-4" />
          </Link>
        </div>
      )}

      {/* Stats Grid */}
      <div className="grid gap-lg md:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => {
          const Icon = stat.icon;
          return (
            <Link key={stat.label} to={stat.href}>
              <Card className="hover:shadow-card-hover cursor-pointer">
                <CardContent className="p-lg">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-caption text-text-secondary uppercase tracking-wide">{stat.label}</p>
                      <p className="text-title2 font-bold text-text-primary mt-xs">{stat.value}</p>
                    </div>
                    <div className={`p-md rounded-lg bg-surface-secondary ${stat.color}`}>
                      <Icon className="h-5 w-5" />
                    </div>
                  </div>
                </CardContent>
              </Card>
            </Link>
          );
        })}
      </div>

      {/* Status Cards */}
      <div className="grid gap-lg md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Request Status</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-md">
              <div className="flex items-center justify-between p-md rounded-lg bg-surface-secondary">
                <div className="flex items-center gap-md">
                  <CheckCircle className="h-5 w-5 text-success" />
                  <span className="text-body text-text-secondary">Completed</span>
                </div>
                <span className="text-headline font-semibold text-success">{completedRequests}</span>
              </div>
              <div className="flex items-center justify-between p-md rounded-lg bg-surface-secondary">
                <div className="flex items-center gap-md">
                  <XCircle className="h-5 w-5 text-error" />
                  <span className="text-body text-text-secondary">Failed</span>
                </div>
                <span className="text-headline font-semibold text-error">{failedRequests}</span>
              </div>
              <div className="flex items-center justify-between p-md rounded-lg bg-surface-secondary">
                <div className="flex items-center gap-md">
                  <Ban className="h-5 w-5 text-warning" />
                  <span className="text-body text-text-secondary">Cancelled</span>
                </div>
                <span className="text-headline font-semibold text-warning">{cancelledRequests}</span>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Quick Actions</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-md">
              <Link
                to="/providers"
                className="flex items-center justify-between p-md rounded-lg bg-surface-secondary hover:bg-surface-hover transition-colors group"
              >
                <div className="flex items-center gap-md">
                  <Server className="h-5 w-5 text-info" />
                  <span className="text-body text-text-primary">Manage Providers</span>
                </div>
                <ArrowRight className="h-4 w-4 text-text-muted group-hover:text-text-primary transition-colors" />
              </Link>
              <Link
                to="/routes"
                className="flex items-center justify-between p-md rounded-lg bg-surface-secondary hover:bg-surface-hover transition-colors group"
              >
                <div className="flex items-center gap-md">
                  <Route className="h-5 w-5 text-accent" />
                  <span className="text-body text-text-primary">Configure Routes</span>
                </div>
                <ArrowRight className="h-4 w-4 text-text-muted group-hover:text-text-primary transition-colors" />
              </Link>
              <Link
                to="/requests"
                className="flex items-center justify-between p-md rounded-lg bg-surface-secondary hover:bg-surface-hover transition-colors group"
              >
                <div className="flex items-center gap-md">
                  <Activity className="h-5 w-5 text-success" />
                  <span className="text-body text-text-primary">View Requests</span>
                </div>
                <ArrowRight className="h-4 w-4 text-text-muted group-hover:text-text-primary transition-colors" />
              </Link>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Features */}
      {!hasProviders && (
        <div className="grid grid-cols-3 gap-lg">
          <div className="bg-surface-secondary/50 rounded-xl p-lg text-center">
            <div className="w-10 h-10 rounded-lg bg-success/10 flex items-center justify-center mx-auto mb-md">
              <CheckCircle className="h-5 w-5 text-success" />
            </div>
            <h3 className="text-headline font-semibold text-text-primary">Secure</h3>
            <p className="text-caption text-text-secondary mt-xs">End-to-end encryption</p>
          </div>
          <div className="bg-surface-secondary/50 rounded-xl p-lg text-center">
            <div className="w-10 h-10 rounded-lg bg-accent/10 flex items-center justify-center mx-auto mb-md">
              <Zap className="h-5 w-5 text-accent" />
            </div>
            <h3 className="text-headline font-semibold text-text-primary">Fast</h3>
            <p className="text-caption text-text-secondary mt-xs">Low latency routing</p>
          </div>
          <div className="bg-surface-secondary/50 rounded-xl p-lg text-center">
            <div className="w-10 h-10 rounded-lg bg-info/10 flex items-center justify-center mx-auto mb-md">
              <Activity className="h-5 w-5 text-info" />
            </div>
            <h3 className="text-headline font-semibold text-text-primary">Insights</h3>
            <p className="text-caption text-text-secondary mt-xs">Real-time analytics</p>
          </div>
        </div>
      )}
        </div>
      </div>
    </div>
  );
}
