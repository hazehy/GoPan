import json
import sys
import os
import matplotlib.pyplot as plt

def safe_get(d, *keys, default=None):
    for k in keys:
        if isinstance(d, dict) and k in d:
            d = d[k]
        else:
            return default
    return d

def main():
    if len(sys.argv) < 2:
        print('Usage: plot_k6_summary.py <k6_summary.json>')
        sys.exit(2)

    path = sys.argv[1]
    with open(path, 'r', encoding='utf-8') as f:
        data = json.load(f)

    metrics = data.get('metrics', data)

    total_reqs = safe_get(metrics, 'http_reqs', 'values', 'count', default=0)
    req_rate = safe_get(metrics, 'http_reqs', 'values', 'rate', default=0)

    duration_vals = safe_get(metrics, 'http_req_duration', 'values', default={}) or {}
    avg = duration_vals.get('avg') or duration_vals.get('mean')
    med = duration_vals.get('med') or duration_vals.get('p(50)')
    p95 = duration_vals.get('p(95)')
    p99 = duration_vals.get('p(99)')

    failed_rate = safe_get(metrics, 'http_req_failed', 'values', 'rate', default=0)

    # Plot counts and rates
    fig, axes = plt.subplots(1, 2, figsize=(12, 5))

    axes[0].bar(['total_requests', 'req_rate', 'failed_rate'], [total_reqs, req_rate, failed_rate], color=['#4C72B0', '#55A868', '#C44E52'])
    axes[0].set_title('Requests / Rate / Failed Rate')

    # Latency chart (ms)
    latency_keys = []
    latency_vals = []
    if avg is not None:
        latency_keys.append('avg')
        latency_vals.append(avg)
    if med is not None:
        latency_keys.append('median')
        latency_vals.append(med)
    if p95 is not None:
        latency_keys.append('p95')
        latency_vals.append(p95)
    if p99 is not None:
        latency_keys.append('p99')
        latency_vals.append(p99)

    if latency_keys:
        axes[1].bar(latency_keys, [v * 1000 for v in latency_vals], color='#8172B2')
        axes[1].set_title('Latency (ms)')
    else:
        axes[1].text(0.5, 0.5, 'No latency metrics found', ha='center')

    plt.suptitle(f'k6 Summary: {os.path.basename(path)}')

    out_png = os.path.splitext(path)[0] + '.png'
    plt.tight_layout(rect=[0, 0.03, 1, 0.95])
    plt.savefig(out_png)
    print('Saved chart to', out_png)

if __name__ == '__main__':
    main()
