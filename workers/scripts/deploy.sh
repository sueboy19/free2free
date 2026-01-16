#!/bin/bash

echo "🚀 Deploying Free2Free API to Cloudflare Workers..."

# 执行 TypeScript 编译
echo "📦 Building TypeScript..."
npm run typecheck

if [ $? -ne 0 ]; then
  echo "❌ TypeScript compilation failed"
  exit 1
fi

# 运行测试
echo "🧪 Running tests..."
npm test

if [ $? -ne 0 ]; then
  echo "❌ Tests failed"
  exit 1
fi

# 运行 lint
echo "🔍 Running linter..."
npm run lint

if [ $? -ne 0 ]; then
  echo "❌ Linting failed"
  exit 1
fi

# 部署到 Cloudflare Workers
echo "🌍 Deploying to Cloudflare Workers..."
wrangler deploy

if [ $? -eq 0 ]; then
  echo "✅ Deployment successful!"
else
  echo "❌ Deployment failed"
  exit 1
fi
