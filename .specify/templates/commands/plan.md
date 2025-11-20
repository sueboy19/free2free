# speckit.plan Command Execution Workflow

Description: Execute the implementation planning workflow using the plan template to generate design artifacts.

## User Input

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

## Outline

1. **Setup**: Run `.specify/scripts/powershell/setup-plan.ps1 -Json` from repo root and parse JSON for FEATURE_SPEC, IMPL_PLAN, SPECS_DIR, BRANCH. For single quotes in args like "I'm Groot", use escape syntax: e.g 'I'\''m Groot' (or double-quote if possible: "I'm Groot").

2. **Load context**: Read FEATURE_SPEC and `.specify/memory/constitution.md`. Load IMPL_PLAN template (already copied).

3. **Execute plan workflow**: Follow the structure in IMPL_PLAN template to:
   - Fill Technical Context (mark unknowns as "NEEDS CLARIFICATION")
   - Fill Constitution Check section from constitution
   - Evaluate gates (ERROR if violations unjustified)
   - Phase 0: Generate research.md (resolve all NEEDS CLARIFICATION)
   - Phase 1: Generate data-model.md, contracts/, quickstart.md
   - Phase 1: Update agent context by running the agent script
   - Re-evaluate Constitution Check post-design

4. **Stop and report**: Command ends after Phase 2 planning. Report branch, IMPL_PLAN path, and generated artifacts.

## Phases

### Phase 0: Outline & Research

1. **Extract unknowns from Technical Context** above:
   - For each NEEDS CLARIFICATION → research task
   - For each dependency → best practices task
   - For each integration → patterns task

2. **Generate and dispatch research agents**:

   ```text
   For each unknown in Technical Context:
     Task: "Research {unknown} for {feature context}"
   For each technology choice:
     Task: "Find best practices for {tech} in {domain}"
   ```

3. **Consolidate findings** in `research.md` using format:
   - Decision: [what was chosen]
   - Rationale: [why chosen]
   - Alternatives considered: [what else evaluated]

**Output**: research.md with all NEEDS CLARIFICATION resolved

### Phase 1: Design & Contracts

**Prerequisites:** `research.md` complete

1. **Extract entities from feature spec** → `data-model.md`:
   - Entity name, fields, relationships
   - Validation rules from requirements
   - State transitions if applicable

2. **Generate API contracts** from functional requirements:
   - For each user action → endpoint
   - Use standard REST/GraphQL patterns
   - Output OpenAPI/GraphQL schema to `/contracts/`

3. **Agent context update**:
   - Run `.specify/scripts/powershell/update-agent-context.ps1 -AgentType roo`
   - These scripts detect which AI agent is in use
   - Update the appropriate agent-specific context file
   - Add only new technology from current plan
   - Preserve manual additions between markers

**Output**: data-model.md, /contracts/*, quickstart.md, agent-specific file

## Constitution Check Guidelines

### Pre-Phase 0 Gate (MUST PASS)

Before starting any research, the implementation plan must pass the Constitution Check:

1. **模組化設計優先** (NON-NEGOTIABLE):
   - Verify the proposed architecture follows modular design principles
   - Ensure code separation into distinct packages
   - Check that no monolithic code structures are planned

2. **API 文件優先** (NON-NEGOTIABLE):
   - Confirm all API endpoints will have Swagger documentation
   - Verify RESTful principles will be followed
   - Check error handling and status codes are planned

3. **測試驅動開發** (NON-NEGOTIABLE):
   - Ensure TDD approach is planned (test-first methodology)
   - Verify comprehensive test coverage strategy
   - Check for unit, integration, and e2e test plans

4. **安全性和認證優先**:
   - Verify authentication and authorization mechanisms are planned
   - Check OAuth 2.0 and JWT implementation strategy
   - Ensure security vulnerabilities will be addressed

5. **可擴展性和性能**:
   - Verify scalability considerations are included
   - Check performance requirements and optimization plans
   - Ensure response time targets are achievable

### Post-Phase 1 Re-evaluation

After completing the design phase, re-evaluate the Constitution Check:

- Update all checkboxes based on final design decisions
- Document any constitution violations and their justifications
- Ensure all NON-NEGOTIABLE principles are fully compliant
- Add any new compliance considerations discovered during design

### Gate Failure Handling

If Constitution Check fails:
- **ERROR** must be thrown if violations are unjustified
- Document all violations in the Complexity Tracking section
- Provide detailed justification for any necessary deviations
- Consider alternative approaches that maintain compliance

## Key rules

- Use absolute paths
- ERROR on gate failures or unresolved clarifications
- Constitution Check must be completed twice: before Phase 0 and after Phase 1
- All NON-NEGOTIABLE principles must be fully compliant
- Document any deviations with clear justification