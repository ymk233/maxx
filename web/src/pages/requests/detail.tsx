import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Button, Card, CardContent, CardHeader, CardTitle, Badge } from '@/components/ui';
import { useProxyRequest, useProxyUpstreamAttempts, useProxyRequestUpdates } from '@/hooks/queries';
import { ArrowLeft, Clock, Zap, FileText, AlertCircle, Server, RotateCcw } from 'lucide-react';
import { statusVariant } from './index';
import type { ProxyUpstreamAttempt } from '@/lib/transport';

export function RequestDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { data: request, isLoading, error } = useProxyRequest(Number(id));
  const { data: attempts } = useProxyUpstreamAttempts(Number(id));
  const [selectedAttemptId, setSelectedAttemptId] = useState<number | null>(null);

  // Subscribe to real-time updates
  useProxyRequestUpdates();

  // 获取选中的 attempt
  const selectedAttempt = attempts?.find(a => a.id === selectedAttemptId) ?? attempts?.[0];

  const formatDuration = (ns: number) => {
    const ms = ns / 1_000_000;
    if (ms < 1000) return `${ms.toFixed(0)}ms`;
    return `${(ms / 1000).toFixed(2)}s`;
  };

  const formatDateTime = (dateStr: string) => {
    if (!dateStr || dateStr === '0001-01-01T00:00:00Z') return '-';
    const date = new Date(dateStr);
    return date.toLocaleString();
  };

  const formatJSON = (obj: unknown): string => {
    if (!obj) return '-';
    try {
      return JSON.stringify(obj, null, 2);
    } catch {
      return String(obj);
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center p-12">
        <p className="text-gray-500">Loading...</p>
      </div>
    );
  }

  if (error || !request) {
    return (
      <div className="space-y-6">
        <Button variant="outline" onClick={() => navigate('/requests')}>
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to Requests
        </Button>
        <Card>
          <CardContent className="p-12 text-center">
            <AlertCircle className="mx-auto mb-4 h-12 w-12 text-red-500" />
            <h3 className="text-lg font-semibold">Request Not Found</h3>
            <p className="text-gray-500">The request you're looking for doesn't exist or has been deleted.</p>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Button variant="outline" onClick={() => navigate('/requests')}>
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back
        </Button>
        <h2 className="text-2xl font-bold">Request Details</h2>
        <Badge variant={statusVariant[request.status]}>{request.status}</Badge>
      </div>

      {/* Overview */}
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 text-sm text-gray-500">
              <Clock className="h-4 w-4" />
              Duration
            </div>
            <p className="mt-1 text-2xl font-semibold">
              {request.duration ? formatDuration(request.duration) : '-'}
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 text-sm text-gray-500">
              <Zap className="h-4 w-4" />
              Tokens
            </div>
            <p className="mt-1 text-2xl font-semibold">
              {request.inputTokenCount + request.outputTokenCount > 0
                ? `${request.inputTokenCount} / ${request.outputTokenCount}`
                : '-'}
            </p>
            <p className="text-xs text-gray-400">Input / Output</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 text-sm text-gray-500">
              <FileText className="h-4 w-4" />
              Cache
            </div>
            <p className="mt-1 text-2xl font-semibold">
              {request.cacheReadCount + request.cacheWriteCount > 0
                ? `${request.cacheReadCount} / ${request.cacheWriteCount}`
                : '-'}
            </p>
            <p className="text-xs text-gray-400">Read / Write</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="text-sm text-gray-500">Cost</div>
            <p className="mt-1 text-2xl font-semibold">
              {request.cost > 0 ? `$${request.cost.toFixed(4)}` : '-'}
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Basic Info */}
      <Card>
        <CardHeader>
          <CardTitle>Basic Information</CardTitle>
        </CardHeader>
        <CardContent>
          <dl className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            <div>
              <dt className="text-sm font-medium text-gray-500">ID</dt>
              <dd className="mt-1 font-mono text-sm">{request.id}</dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-gray-500">Request ID</dt>
              <dd className="mt-1 font-mono text-xs break-all">{request.requestID || '-'}</dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-gray-500">Session ID</dt>
              <dd className="mt-1 font-mono text-xs break-all">{request.sessionID || '-'}</dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-gray-500">Client Type</dt>
              <dd className="mt-1">
                <Badge variant="info">{request.clientType}</Badge>
              </dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-gray-500">Request Model</dt>
              <dd className="mt-1 font-mono text-sm">{request.requestModel || '-'}</dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-gray-500">Response Model</dt>
              <dd className="mt-1 font-mono text-sm">{request.responseModel || '-'}</dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-gray-500">Start Time</dt>
              <dd className="mt-1 text-sm">{formatDateTime(request.startTime)}</dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-gray-500">End Time</dt>
              <dd className="mt-1 text-sm">{formatDateTime(request.endTime)}</dd>
            </div>
            <div>
              <dt className="text-sm font-medium text-gray-500">Retry Attempts</dt>
              <dd className="mt-1 text-sm">{request.proxyUpstreamAttemptCount}</dd>
            </div>
          </dl>
        </CardContent>
      </Card>

      {/* Error */}
      {request.error && (
        <Card className="border-red-200 bg-red-50 dark:border-red-900 dark:bg-red-950">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-red-600 dark:text-red-400">
              <AlertCircle className="h-5 w-5" />
              Error
            </CardTitle>
          </CardHeader>
          <CardContent>
            <pre className="whitespace-pre-wrap break-words font-mono text-sm text-red-600 dark:text-red-400">
              {request.error}
            </pre>
          </CardContent>
        </Card>
      )}

      {/* Upstream Attempts - 左右布局 */}
      {attempts && attempts.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <RotateCcw className="h-5 w-5" />
              Upstream Attempts ({attempts.length})
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex gap-4">
              {/* 左侧：Attempts 列表 */}
              <div className="w-64 shrink-0 space-y-2">
                {attempts.map((attempt: ProxyUpstreamAttempt, index: number) => (
                  <div
                    key={attempt.id}
                    onClick={() => setSelectedAttemptId(attempt.id)}
                    className={`cursor-pointer rounded-lg border p-3 transition-colors hover:bg-gray-50 dark:hover:bg-gray-800 ${
                      (selectedAttempt?.id === attempt.id) ? 'border-blue-500 bg-blue-50 dark:bg-blue-950' : ''
                    }`}
                  >
                    <div className="flex items-center gap-2">
                      <span className="flex h-6 w-6 items-center justify-center rounded-full bg-gray-100 text-sm font-medium dark:bg-gray-700">
                        {index + 1}
                      </span>
                      <Badge variant={statusVariant[attempt.status as keyof typeof statusVariant] || 'secondary'}>
                        {attempt.status}
                      </Badge>
                    </div>
                    <div className="mt-2 flex items-center gap-1 text-xs text-gray-500">
                      <Server className="h-3 w-3" />
                      Route #{attempt.routeID} / Provider #{attempt.providerID}
                    </div>
                    {/* Token Stats */}
                    {(attempt.inputTokenCount > 0 || attempt.outputTokenCount > 0) && (
                      <div className="mt-1 text-xs text-gray-400">
                        {attempt.inputTokenCount} / {attempt.outputTokenCount} tokens
                      </div>
                    )}
                  </div>
                ))}
              </div>

              {/* 右侧：选中 Attempt 的详情 */}
              <div className="min-w-0 flex-1 space-y-4">
                {selectedAttempt ? (
                  <>
                    {/* Request Info */}
                    {selectedAttempt.requestInfo && (
                      <div className="rounded-lg border p-4">
                        <h4 className="mb-3 font-medium">Request</h4>
                        <div className="mb-3 flex items-center gap-2">
                          <Badge variant="info">{selectedAttempt.requestInfo.method}</Badge>
                          <code className="flex-1 truncate rounded bg-gray-100 px-2 py-1 font-mono text-xs dark:bg-gray-800">
                            {selectedAttempt.requestInfo.url}
                          </code>
                        </div>
                        <div className="mb-2">
                          <h5 className="mb-1 text-sm text-gray-500">Headers</h5>
                          <pre className="max-h-32 overflow-auto rounded bg-gray-100 p-2 font-mono text-xs dark:bg-gray-800">
                            {formatJSON(selectedAttempt.requestInfo.headers)}
                          </pre>
                        </div>
                        {selectedAttempt.requestInfo.body && (
                          <div>
                            <h5 className="mb-1 text-sm text-gray-500">Body</h5>
                            <pre className="max-h-48 overflow-auto rounded bg-gray-100 p-2 font-mono text-xs dark:bg-gray-800">
                              {(() => {
                                try {
                                  return formatJSON(JSON.parse(selectedAttempt.requestInfo.body));
                                } catch {
                                  return selectedAttempt.requestInfo.body;
                                }
                              })()}
                            </pre>
                          </div>
                        )}
                      </div>
                    )}

                    {/* Response Info */}
                    {selectedAttempt.responseInfo && (
                      <div className="rounded-lg border p-4">
                        <h4 className="mb-3 font-medium">Response</h4>
                        <div className="mb-3">
                          <Badge variant={selectedAttempt.responseInfo.status >= 400 ? 'danger' : 'success'}>
                            {selectedAttempt.responseInfo.status}
                          </Badge>
                        </div>
                        <div className="mb-2">
                          <h5 className="mb-1 text-sm text-gray-500">Headers</h5>
                          <pre className="max-h-32 overflow-auto rounded bg-gray-100 p-2 font-mono text-xs dark:bg-gray-800">
                            {formatJSON(selectedAttempt.responseInfo.headers)}
                          </pre>
                        </div>
                        {selectedAttempt.responseInfo.body && (
                          <div>
                            <h5 className="mb-1 text-sm text-gray-500">Body</h5>
                            <pre className="max-h-48 overflow-auto rounded bg-gray-100 p-2 font-mono text-xs dark:bg-gray-800">
                              {(() => {
                                try {
                                  return formatJSON(JSON.parse(selectedAttempt.responseInfo.body));
                                } catch {
                                  return selectedAttempt.responseInfo.body;
                                }
                              })()}
                            </pre>
                          </div>
                        )}
                      </div>
                    )}

                    {/* 没有请求/响应信息 */}
                    {!selectedAttempt.requestInfo && !selectedAttempt.responseInfo && (
                      <div className="flex items-center justify-center rounded-lg border border-dashed p-12 text-gray-400">
                        No request/response data available
                      </div>
                    )}
                  </>
                ) : (
                  <div className="flex items-center justify-center rounded-lg border border-dashed p-12 text-gray-400">
                    Select an attempt to view details
                  </div>
                )}
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Request Info */}
      {request.requestInfo && (
        <Card>
          <CardHeader>
            <CardTitle>Request</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center gap-4">
              <Badge variant="info">{request.requestInfo.method}</Badge>
              <code className="flex-1 truncate rounded bg-gray-100 px-2 py-1 font-mono text-sm dark:bg-gray-800">
                {request.requestInfo.url}
              </code>
            </div>
            <div>
              <h4 className="mb-2 font-medium">Headers</h4>
              <pre className="max-h-48 overflow-auto rounded bg-gray-100 p-3 font-mono text-xs dark:bg-gray-800">
                {formatJSON(request.requestInfo.headers)}
              </pre>
            </div>
            {request.requestInfo.body && (
              <div>
                <h4 className="mb-2 font-medium">Body</h4>
                <pre className="max-h-96 overflow-auto rounded bg-gray-100 p-3 font-mono text-xs dark:bg-gray-800">
                  {(() => {
                    try {
                      return formatJSON(JSON.parse(request.requestInfo.body));
                    } catch {
                      return request.requestInfo.body;
                    }
                  })()}
                </pre>
              </div>
            )}
          </CardContent>
        </Card>
      )}

      {/* Response Info */}
      {request.responseInfo && (
        <Card>
          <CardHeader>
            <CardTitle>Response</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center gap-4">
              <Badge variant={request.responseInfo.status >= 400 ? 'danger' : 'success'}>
                {request.responseInfo.status}
              </Badge>
            </div>
            <div>
              <h4 className="mb-2 font-medium">Headers</h4>
              <pre className="max-h-48 overflow-auto rounded bg-gray-100 p-3 font-mono text-xs dark:bg-gray-800">
                {formatJSON(request.responseInfo.headers)}
              </pre>
            </div>
            {request.responseInfo.body && (
              <div>
                <h4 className="mb-2 font-medium">Body</h4>
                <pre className="max-h-96 overflow-auto rounded bg-gray-100 p-3 font-mono text-xs dark:bg-gray-800">
                  {(() => {
                    try {
                      return formatJSON(JSON.parse(request.responseInfo.body));
                    } catch {
                      return request.responseInfo.body;
                    }
                  })()}
                </pre>
              </div>
            )}
          </CardContent>
        </Card>
      )}
    </div>
  );
}
